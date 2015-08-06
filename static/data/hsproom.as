#ifndef hsproom
#define hsproom

#module hsproom

#deffunc meslog str m
	mestmp = m
	bsave "/dev/stdout", mestmp, strlen(mestmp)
return

#deffunc meserr str m
	mestmp = m
	bsave "/dev/stderr", mestmp, strlen(mestmp)
return

#deffunc prompt var buf
	buf = ""

	repeat
		tmp = ""

		bload "/dev/stdin", tmp, 1

		if (peek(tmp,0) == 10) {
			break
		}

		buf += tmp
	loop

	return tmp
return

#global

#endif
