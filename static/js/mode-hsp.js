ace.define("ace/mode/hsp_highlight_rules",["require","exports","module","ace/lib/oop","ace/mode/text_highlight_rules"], function(require, exports, module) {
  "use strict";

  var oop = require("../lib/oop");
  var TextHighlightRules = require("./text_highlight_rules").TextHighlightRules;

  var HspHighlightRules = function() {

	this.$rules = {
      "start": [
		{
          token: [
			"meta.ending-space"
          ],
          regex: "$"
		},
        {
          token: "comment",
          regex: "(;|//).*$",
        },
        {
          token : "comment", // multi line comment
          regex : "\\/\\*",
          next : "comment"
        },
		{
          token : "keyword", // pre-compiler directives
		  regex : "#\\s*(?:define|func|cfunc|include|addition|uselib|global|module|deffunc|defcfunc|pack|epack|packopt|const|undef|if|ifdef|ifndef|else|endif|modfunc|modcfunc|modinit|modterm|regcmd|cmd|usecom|comfunc|enum|runtime|cmpopt|aht|ahtmes)\\b",
          next  : "directive"
		},
		{
          token : "constant.other", // label
          regex : "\\*(\\b[a-zA-Z_x7f-xff][a-zA-Z0-9_x7f-xff]*?\\b(?=\\s*((;|\\/\\/|\\/\\*).*)?$))"
		},
        {
          token: [
            "keyword.control"
          ],
		  regex: "\\b(?:await|exec|end|gosub|goto|if|loop|onexit|return|run|stop|wait|repeat|break|continue|onkey|onclick|onerror|exgoto|on|foreach|oncmd|else|while|wend|until|do|for|next|switch|swend|default|case|swbreak|_continue|_break)\\b"
        },
        {
          token: "punctuation.definition.string.begin",
          regex: '"',
          next: "string"
        },
        {
          token: "support.function",
          regex: "\\b(?:comres|alloc|dim|dimtype|poke|wpoke|lpoke|sdim|ddim|memcpy|memset|newmod|delmod|memexpand|ldim|newlab|bcopy|chdir|delete|dirlist|exist|mkdir|bload|bsave|memfile|chdpm|cls|mes|print|title|dialog|bgscr|bmpsave|boxf|buffer|chgdisp|color|font|gcopy|gmode|gsel|gzoom|palcolor|palette|pget|picload|pos|pset|redraw|screen|width|sysfont|line|circle|syscolor|hsvcolor|grect|grotate|gsquare|axobj|winobj|sendmsg|groll|gradf|celload|celdiv|celput|ginfo|objinfo|int|rnd|strlen|length|length2|length3|length4|vartype|varptr|gettime|str|dirinfo|double|sin|cos|tan|atan|sqrt|sysinfo|callfunc|absf|abs|logf|expf|limit|limitf|varuse|libptr|powf|getease|geteasef|dup|dupptr|mref|mci|mmplay|mmload|mmstop|button|chkbox|clrobj|combox|input|listbox|mesbox|objprm|objsize|objsel|objmode|objenable|objskip|objimage|peek|wpeek|lpeek|getkey|mouse|randomize|stick|mcall|setease|logmes|assert|getstr|noteadd|notedel|noteget|notesel|noteunsel|notesave|noteload|cnvstow|split|strrep|noteinfo|strmid|instr|getpath|strf|cnvwtos|strtrim|setreq|getreq|gfilter|mtlist|mtinfo|devinfo|devinfoi|devprm|devcontrol|gettimestr|getdatestr|note2array|array2note|arraysave|arrayload|emath|emstr|emcnv|emint|emsin|emcos|emsqr|ematan|pipeexec|pipeget|pipeput|mmvol|mmpan|mmstat|setcls|celputm|fvseti|fvset|fvadd|fvsub|fvmul|fvdiv|fvdir|fvface|fvmin|fvmax|fvouter|fvinner|fvunit|fsin|fcos|fsqr|str2fv|fv2str|str2f|f2str|delobj|setpos|setang|setangr|setscale|setdir|setwork|addpos|addang|addangr|addscale|adddir|addwork|getpos|getscale|getdir|getwork|getposi|getscalei|getdiri|getworki|selpos|selang|selscale|seldir|selwork|objset3|objsetf3|objadd3|objaddf3|objadd3r|objset3r|setobjmode|setcoli|getcoli|getobjcoli|findobj|nextobj|setborder|selmoc|objgetfv|objsetfv|objaddfv|objexist|gpreset|gpdraw|gpusescene|gpsetprm|gpgetprm|gppostefx|gpuselight|gpusecamera|gpmatprm|gpmatstate|gpviewport|setobjname|getobjname|gpcolormat|gptexmat|gpusermat|gpclone|gpnull|gpload|gpplate|gpfloor|gpbox|gpspr|gplight|gpcamera|gplookat|gppbind|getwork2|getwork2i|selquat|selwork2|setwork2|addwork2|gpcnvaxis|gppset|gpobjpool|gppapply|gpmatprm1|gpmatprm4|setalpha|bmppalette|text|textmode|emes|gfade|statictext|statictext_set|scrollbar|progbar|progbar_set|progbar_step|gfini|gfcopy|gfdec|gfinc|getobjsize|resizeobj|objgray|p_scrwnd)\\b"
        },
        {
          token: [
            "constant.numeric"
          ],
          regex: "-?\\b(?:(?:(?:0(?:x|X)|(\\\$))[0-9a-fA-F]*)|(?:(?:[0-9]+\\.?[0-9]*)|(?:\\.[0-9]+))(?:(?:e|E)(?:\\+|-)?[0-9]+)?)(?:F|f)?\\b"
        },
        {
          token: [
            "variable"
          ],
          regex: "\\b(?:_debug|__hsp30__|__file__|__line__|__date__|__time__|__hspver__|gmode_sub|gmode_add|gmode_gdi|gmode_rgb0|gmode_mem|gmode_alpha|gmode_rgb0alpha|gmode_pixela|objinfo_mode|objinfo_bmscr|objinfo_hwnd|screen_normal|screen_palette|screen_hide|screen_fixedsize|screen_tool|screen_frame|font_normal|font_bold|font_italic|font_underline|font_strikeout|font_antialias|objmode_normal|objmode_guifont|objmode_usefont|and|or|xor|not|rad2deg|deg2rad|msgothic|msmincho|hspstat|hspver|cnt|err|stat|mousex|mousey|mousew|strsize|refstr|looplev|sublev|iparam|wparam|lparam|hwnd|hdc|hinstance|refdval|thismod|notemax|notesize|ginfo_mx|ginfo_my|ginfo_act|ginfo_sel|ginfo_wx1|ginfo_wy1|ginfo_wx2|ginfo_wy2|ginfo_vx|ginfo_vy|ginfo_sizex|ginfo_sizey|ginfo_winx|ginfo_winy|ginfo_sx|ginfo_sy|ginfo_mesx|ginfo_mesy|ginfo_r|ginfo_g|ginfo_b|ginfo_paluse|ginfo_dispx|ginfo_dispy|ginfo_cx|ginfo_cy|ginfo_intid|ginfo_newid|ginfo_accx|ginfo_accy|ginfo_accz|dir_cur|dir_exe|dir_win|dir_sys|dir_cmdline|dir_desktop|dir_mydoc|dir_tv|M_PI)\\b"
        },
        {
          token: [
            "keyword.operator"
          ],
          regex: "\\-|\\+|\\*|\\/|\\>|\\<|\\=|\\&"
        }
      ],
      "comment": [
        {
          token : "comment", // closing comment
          regex : ".*?\\*\\/",
          next : "start"
        }, {
          token : "comment", // comment spanning whole line
          regex : ".+"
        }
      ],
      "string": [
        {
          token: "constant.character.escape.apostrophe",
          regex: '""'
        },
        {
          token: "string.quoted.double",
          regex: '"',
          next: "start"
        },
        {
          defaultToken: "string.quoted.double"
        }
      ],
      "directive" : [
        {
          token : "constant.other.multiline",
          regex : /\\/
        },
        {
          token : "constant.other.multiline",
          regex : /.*\\/
        },
            {
              token : "constant.other",
              regex : "\\s*<.+?>",
              next : "start"
            },
        {
          token : "constant.other", // single line
          regex : '\\s*["](?:(?:\\\\.)|(?:[^"\\\\]))*?["]',
          next : "start"
        }, 
        {
                token : "constant.other", // single line
          regex : "\\s*['](?:(?:\\\\.)|(?:[^'\\\\]))*?[']",
          next : "start"
        },
        {
          token : "constant.other",
          regex : /[^\\\/]+/,
          next : "start"
        }
      ]
    }

  };

  oop.inherits(HspHighlightRules, TextHighlightRules);

  exports.HspHighlightRules = HspHighlightRules;
});

ace.define("ace/mode/hsp",["require","exports","module","ace/lib/oop","ace/mode/text","ace/mode/hsp_highlight_rules"], function(require, exports, module) {
"use strict";

var oop = require("../lib/oop");
var TextMode = require("./text").Mode;
var HspHighlightRules = require("./hsp_highlight_rules").HspHighlightRules;

var Mode = function() {
  this.HighlightRules = HspHighlightRules;
};
oop.inherits(Mode, TextMode);

(function() {
  this.lineCommentStart = [";", "//"];
  this.blockComment = {start: "/*", end: "*/"};

  this.$id = "ace/mode/hsp";
}).call(Mode.prototype);

exports.Mode = Mode;
});
