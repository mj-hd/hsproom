(function($) {

	$.fn.uploadbox = function(options){
		 
		var defaults = {
			allowMultiFiles: false,
			emptyText: "ファイルをドロップしてください。",
			emptyIconClass:   "mdi-editor-attach-file",
			inputButtonClass: "mdi-file-file-upload",
			limitSize: 1024*1024*5
		};
		var setting = $.extend(defaults, options);

		this.each(function() {
			var $this = $(this);
			$this.empty();
			$this.addClass("upload_frame").addClass("well");
			/*$this.append(
				$("<div />").addClass("progress").addClass("progress-striped").addClass("active").append(
					$("<div />").addClass("upload_progress").addClass("progress-bar").css({
						"width": "0%"
					})
					)
				);*/
			$this.append(
				$("<div />").addClass("upload_box").append(
					$("<div />").addClass("upload_box_nofile").append(
						$("<i />").addClass(setting.emptyIconClass)
						).append(
						$("<br />")
						)
					).append(
						setting.allowMultiFiles ? $('<ul style="display:none;"/>').addClass("list-unstyled") : $('<img width="100%" style="display: none" />').addClass("upload_box_preview")
					).append(
						$("<span />").addClass("upload_box_status").text(setting.emptyText)
						)
				);
			$this.append(
				$("<span />").addClass("btn").addClass("btn-flat").addClass("btn-default").addClass("upload_box_fileinput_button").append(
					$("<i />").addClass(setting.inputButtonClass)
					).append(
						setting.allowMultiFiles ? $('<input type="file" multiple />') : $('<input type="file" />')
					)
				);

			if (setting.allowMultiFiles) {
				$this.data("totalSize", 0);
				$this.data("multiFiles", true);
				$this.data("fileCount", 0);
			} else {
				$this.data("multiFiles", false);
			}

			// D&D
			$this.on("drop", function(e) {
				e.preventDefault();

				var files = e.originalEvent.dataTransfer.files;

				var $$this = $(this);

				$.each(files, function() {
					$$this.trigger("ondrop", this);
				});
			});

			// Input Box
			$this.find("input[type=file]").on("change", function(e) {
				var $$this = $(this);

				$.each(this.files, function() {
					$$this.parents("div .upload_frame").trigger("ondrop", this);
				});
			});

			$this.on("dragover", function(e) {
				e.preventDefault();
			});

			$this.on("ondrop", function(e, file) {

				var box = $this.find(".upload_box");
				var stat = $this.find(".upload_box .upload_box_status");
				var del = $('<a class="badge" href="javascript:void();">X</a>');

				del.on("click", function() {
					$this.resetUploadbox();
				});

				stat.text(file.name + "を読み込み中です...");

				if (file.size > setting.limitSize) {
					stat.text("ファイルサイズが大きすぎます。");
					return;
				}

				if (setting.allowMultiFiles) {
					if ($this.find('ul li[data-fname="' + file.name + '"]').length != 0) {
						stat.text("同じ名前のファイルは登録できません。");
						return;
					}

					var totalSize = $this.data("totalSize");

					if (totalSize + file.size > setting.limitSize) {
						stat.text("ファイルサイズが大きすぎます。");
						return
					}

					$this.data("totalSize", totalSize + file.size);
					del.unbind();
					del.on("click", function() {
						$this.deleteItemUploadbox(file.name);
					});
				}

				var readerUpload = new FileReader();
				readerUpload.onload = function(e) {
					$this.trigger("onupload", [convertArrayBufferToBase64(e.target.result), file]);
					function convertArrayBufferToBase64(buffer) {
						var result = '';
						var bytes = new Uint8Array(buffer);
						var len = bytes.byteLength;
						for (var i = 0; i < len; i++) {
							result += String.fromCharCode(bytes[i]);
						}
						return window.btoa(result);
					}
				};
				readerUpload.readAsArrayBuffer(file);

				if (file.type.indexOf("image") != 0) {
					if (setting.allowMultiFiles) {
						var item = $('<li data-fname="' + file.name + '" data-size="' + file.size + '"/>').append(
									file.name
								).append(
									del
								);

						$this.find("ul").append(item).show();
						$this.find("upload_box_nofile").hide();

						stat.text(setting.emptyText);

						var count = $this.data("fileCount");
						$this.data("fileCount", count+1);
					} else {
						stat.text(file.name);
						stat.append(del);
					}
					return;
				}

				// 画像のプレビュー
				readerPreview = new FileReader();
				readerPreview.onload = function(e) {
					if (setting.allowMultiFiles) {
						var item = $('<li data-fname="' + file.name + '" />').append(
								$('<img height="16px" />').attr("src", e.target.result)
							).append(
								file.name
							).append(
								del
							);

						$this.find("ul").append(item).show();

						stat.text(setting.emptyText);

						$this.find(".upload_box_nofile").hide();

						var count = $this.data("fileCount");
						$this.data("fileCount", count+1);
					} else {
						var preview = $this.find(".upload_box_preview");

						stat.text(file.name);
						stat.append(del);

						preview.attr("src", e.target.result);

						box.find(".upload_box_nofile").hide();
						$this.find(".upload_fileinput_button").hide();

						preview.show();
					}
				};
				readerPreview.readAsDataURL(file);

			});
		});
		 
		return(this);
	};

	$.fn.resetUploadbox = function() {
		this.find(".upload_box_status").text("ファイルをドロップしてください。");
		this.find(".upload_box_fileinput_button").show();
		this.find(".upload_box_nofile").show();

		this.find(".upload_box_preview").attr("src", "").hide();

		return(this);
	};

	$.fn.deleteItemUploadbox = function(fname) {
		var count = this.data("fileCount");

		if (count == 0) {
			return(this);
		}

		var item = this.find('li[data-fname="' + fname  + '"]');
		var fileSize = item.data("size");
		item.remove();

		count--;
		this.data("fileCount", count);

		if (count == 0) {
			this.find(".upload_box_nofile").show();
			this.data("totalSize", 0);
		} else {
			var totalSize = this.data("totalSize");
			this.data("totalSize", totalSize - fileSize);
		}

		this.trigger("ondelete", fname);

		return(this);
	};
})(jQuery);
