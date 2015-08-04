var NTWFS = $.extend(true, {}, MEMFS);

NTWFS.pid = 0;

NTWFS.mount = function (mount) {
	return NTWFS.createNode(null, '/', 16384 | 0777, 0);
}

NTWFS.node_ops.mknod = function (parent, name, mode, dev) {
	return NTWFS.createNode(parent, name, mode, dev);
};

NTWFS.createNode = function (parent, name, mode, dev) {
	var node = FS.createNode(parent, name, mode, dev);

	node.node_ops = NTWFS.node_ops;
	node.stream_ops = NTWFS.stream_ops;
	node.contents = [];
	node.timestamp = Date.now();

	if (parent) parent.contents[name] = node;

	return node;
}

NTWFS.syncfs = function (mount, populate, callback) {
	console.log("syncfs!!");

	console.log("call getLocalSet");
	console.log(FS.readdir("/share"));
	NTWFS.getLocalSet(mount, function (err, local) {
		if (err) return callback(err);

		console.log("call getRemoteSet");
		NTWFS.getRemoteSet(mount, function (err, remote) {
			if (err) return callback(err);

			var src = populate ? remote : local;
			var dst = populate ? local : remote;
	
			NTWFS.reconcile(src, dst, callback);
		});
	});
}

NTWFS.getLocalSet = function (mount, callback) {
	console.log("called getLocalSet");

	var result = { type: "local", entries: [] };

	var dirlist = FS.readdir(mount.mountpoint).filter(isRealDir).map(toAbsolute(mount.mountpoint));

	while (dirlist.length) {
		var path = dirlist.pop();
		var stat;

		try {
			stat = FS.stat(path);
		} catch(e) {
			return callback(e);
		}

		if (FS.isDir(stat.mode)) {
			dirlist.push.apply(dirlist, FS.readdir(path).filter(isRealDir).map(toAbsolute(path)));
		}

		result.entries.push({
			Name: path,
			CreatedAt: stat.ctime,
			UpdatedAt: stat.mtime,
			Mode:      stat.mode
		});
	}

	console.log("getLocalSet returned "+result.entries.length);

	return callback(null, result)

	function isRealDir(p) {
		return p !== '.' && p !== '..';
	};
	function toAbsolute(root) {
		return function(p) {
			return PATH.join2(root, p);
		}
	};
}

NTWFS.getRemoteSet = function (mount, callback) {
	console.log("called getRemoteSet");

	var result = { type: "remote", entries: [] };

	$.ajax({
		method:   "get",
		url:      "/api/program/shared_data/list/",
		dataType: "json",
		data: {
			"p": NTWFS.pid
		}
	}).success(function (data) {
		if (data.DataList) {
			result.entries = data.DataList;

			for (var i = 0; i < result.entries.length; i++) {
				result.entries[i].CreatedAt = new Date(result.entries[i].CreatedAt);
				result.entries[i].UpdatedAt = new Date(result.entries[i].UpdatedAt);
			}
		}
		console.log("getRemoteSet returned "+result.entries.length);
		return callback(null, result);
	}).error(function (data) {
		return callback("failed to download shared file list");
	});
}

NTWFS.reconcile = function (src, dst, callback) {
	console.log(src);
	console.log(dst);

	var create = [];
	var remove = [];
	var total  = 0;

	$.each(src.entries, function(i, sfile) {
		var e;
		$.each(dst.entries, function(j, dfile) {
			if (sfile.Name == dfile.Name) {
				e = dfile;
			}
		});

		if (!e || sfile.UpdatedAt > e.UpdatedAt) {
			create.push(sfile);
			total++;
		}
	});

	$.each(dst.entries, function(i, dfile) {
		var e;
		$.each(src.entries, function(j, sfile) {
			if (dfile.Name == sfile.Name) {
				e = sfile;
			}
		});

		if (!e) {
			remove.push(dfile);
			total++;
		}
	});

	if (total == 0) return callback(null);

	var completed = 0;

	$.each(create.sort(filecomp), function(i, file) {
		if (dst.type === "local") {
			NTWFS.loadRemoteEntry(file, function(err, entry) {
				if (err) return done(err);

				IDBFS.storeLocalEntry(entry.Name, entry, done);
			});
		} else {
			IDBFS.loadLocalEntry(file.Name, function(err, entry) {
				if (err) return done(err);

				console.log(entry.contents);
				NTWFS.storeRemoteEntry(file, entry, done);
			});
		}
	});

	$.each(remove.sort(filecomp), function(i, file) {
		if (dst.type === "local") {
			IDBFS.removeLocalEntry(file.Name, done);
		} else {
			NTWFS.removeRemoteEntry(file, done);
		}
	});

	function done(err) {
		if (err) {
			if (!done.errored) {
				done.errored = true;
				return callback(err);
			}
			return;
		}
		if (++completed >= total) {
			return callback(null);
		}
	}
	function filecomp(a, b) {
		if (a.Name < b.Name) {
			return -1;
		} else if (a.Name > b.Name) {
			return 1;
		}
		return 0;
	}
}

NTWFS.loadRemoteEntry = function (file, callback) {
	$.ajax({
		method: "GET",
		url: "/api/program/shared_data/read/",
		data: {
			"p": NTWFS.pid,
			"f": file.Name
		}
	}).success(function(data) {
		var entry = {};
		entry.contents = data;
		entry.timestamp = file.UpdatedAt;
		entry.mode = file.Mode;
		return callback(null, entry);
	}).error(function(data) {
		return callback("failed to download shared data");
	});
}

NTWFS.storeRemoteEntry = function (file, entry, callback) {
	$.ajax({
		method: "POST",
		url: "/api/program/shared_data/write/",
		dataType: "json",
		data: {
			"p": NTWFS.pid,
			"name": file.Name,
			"data": convertArrayBufferToBase64(entry.contents),
			"created": Math.floor(file.CreatedAt.getTime() / 1000),
			"updated": Math.floor(file.UpdatedAt.getTime() / 1000),
			"mode":    file.Mode
		}
	}).success(function(data) {
		console.log("write entry");
		return callback(null);
	}).error(function(data) {
		return callback("failed to post shared file data");
	});

	function convertArrayBufferToBase64(buffer) {
		var result = '';
		var bytes = new Uint8Array(buffer);
		var len = bytes.byteLength;
		for (var i = 0; i < len; i++) {
			result += String.fromCharCode(bytes[i]);
		}
		return window.btoa(result);
	}
}

NTWFS.removeRemoteEntry = function (file, callback) {
	$.ajax({
		method: "POST",
		url: "/api/program/shared_data/delete/",
		dataType: "json",
		data: {
			"p": NTWFS.pid,
			"f": file.Name
		}
	}).success(function (data) {
		console.log("remove entry");
		return callback(null);
	}).error(function (data) {
		return callback("failed to remove shared file");
	});
}
