/**
 * SyntaxHighlighter
 * http://alexgorbatchev.com/SyntaxHighlighter
 *
 * SyntaxHighlighter is donationware. If you are using it, please donate.
 * http://alexgorbatchev.com/SyntaxHighlighter/donate.html
 *
 * @license
 * Dual licensed under the MIT and GPL licenses.
 */
;(function()
  {
      // CommonJS
      typeof(require) != 'undefined' ? SyntaxHighlighter = require('shCore').SyntaxHighlighter : null;

      function Brush()
      {
	  var r = SyntaxHighlighter.regexLib;
	  var wordlist1 = SyntaxHighlighter.defaults.wordlist1;
	  var wordlist2 = SyntaxHighlighter.defaults.wordlist2;
	  var wordlist3 = SyntaxHighlighter.defaults.wordlist3;

	  this.regexList = [
	      //{ regex: /^\s*@\w+/gm, css: 'decorator' },
	      // userdef
	      { regex: new RegExp(this.getKeywords(wordlist1), 'gm'), css: 'userdef1' },
	      { regex: new RegExp(this.getKeywords(wordlist2), 'gm'), css: 'userdef2' },
	      { regex: new RegExp(this.getKeywords(wordlist3), 'gm'), css: 'userdef3' },
	      // comment
	      { regex: /\/\*[\s\S]*?\*\/|(\/\/|;)[\s\S]*?$/gm, css: 'comments' },
	      { regex: new RegExp('[^&quota|&ampa|&lta|&gta|&OEliga|&oeliga|&Scarona|&scarona|&Yumla|&circa|&tildea|&ndasha|&mdasha|&lsquoa|&rsquoa|&sbquoa|&ldquoa|&rdquoa|&bdquoa|&daggera|&Daggera|&permila|&lsaquoa|&rsaquoa|&euroa|&nbspa|&iexcla|&centa|&pounda|&currena|&yena|&brvbara|&secta|&umla|&copya|&ordfa|&laquoa|&nota|&shya|&rega|&macra|&dega|&plusmna|&sup2a|&sup3a|&acutea|&microa|&paraa|&middota|&cedila|&ordma|&raquoa|&frac14a|&frac12a|&frac34a|&iquesta|&Agravea|&Aacutea|&Acirca|&Atildea|&Aumla|&Aringa|&AEliga|&Ccedila|&Egravea|&Eacutea|&Ecirca|&Eumla|&Igravea|&Iacutea|&Icirca|&Iumla|&ETHa|&Ntildea|&Ogravea|&Oacutea|&Ocirca|&Otildea|&Oumla|&timesa|&Oslasha|&Ugravea|&Uacutea|&Ucirca|&Uumla|&Yacutea|&THORNa|&szliga|&agravea|&aacutea|&acirca|&atildea|&aumla|&aringa|&aeliga|&ccedila|&egravea|&eacutea|&ecirca|&eumla|&igravea|&iacutea|&icirca|&iumla|&etha|&ntildea|&ogravea|&oacutea|&ocirca|&otildea|&oumla|&dividea|&oslasha|&ugravea|&uacutea|&ucirca|&uumla|&yacutea|&thorna|&yumla|&fnofa|&Alphaa|&Betaa|&Gammaa|&Deltaa|&Epsilona|&Zetaa|&Etaa|&Thetaa|&Iotaa|&Kappaa|&Lambdaa|&Mua|&Nua|&Xia|&Omicrona|&Pia|&Rhoa|&Sigmaa|&Taua|&Upsilona|&Phia|&Chia|&Psia|&Omegaa|&alphaa|&betaa|&gammaa|&deltaa|&epsilona|&zetaa|&etaa|&thetaa|&iotaa|&kappaa|&lambdaa|&mua|&nua|&xia|&omicrona|&pia|&rhoa|&sigmafa|&sigmaa|&taua|&upsilona|&phia|&chia|&psia|&omegaa|&bulla|&hellipa|&primea|&Primea|&olinea|&frasla|&tradea|&larra|&uarra|&rarra|&darra|&harra|&rArra|&hArra|&foralla|&parta|&exista|&nablaa|&isina|&nia|&proda|&suma|&minusa|&radica|&propa|&infina|&anga|&anda|&ora|&capa|&cupa|&inta|&there4a|&sima|&asympa|&nea|&equiva|&lea|&gea|&suba|&supa|&subea|&supea|&oplusa|&perpa|&loza|&spadesa|&clubsa|&heartsa|&diamsa];(?!&).+','gm'), css: 'comments' },
	      { regex: r.singleLineCComments, css: 'comments' },
	      { regex: r.multiLineCComments, css: 'comments' },
	      // string
	      { regex: /"(?!")(?:\.|\\\"|[^\""\n])*"/gm, css: 'string' },
	      { regex: /'(?!')(?:\.|(\\\')|[^\''\n])*'/gm, css: 'string' },
	      { regex: /{"([^\\']|\\[\s\S])*"}/g, css: 'string' },
	      // operator
	      { regex: /\+|\-|\*|\/|\\|=|==|\!|&lt;|&gt;|&amp;|\x7c/gm,  css: 'operator' },
	      // number
	      { regex: /\-\$[\da-f]+|\-\b(\d+(\.\d+)?(e(\+|\-)?\d+)?(f|d)?|0x[\da-f]+)\b|\$[\da-f]+|\b(\d+(\.\d+)?(e(\+|\-)?\d+)?(f|d)?|0x[\da-f]+)\b/gi, css: 'value' },
	      // function
	      { regex: /\b(a(bs|bsf|ssert|tan|wait|xobj)|b(copy|gscr|load|mpsave|oxf|reak|save|uffer|utton)|c(allfunc|el(div|load|put)|h(dir|dpm|gdisp|kbox)|ircle|lrobj|ls|n(t|vstow|vwtos)|olor|om(box|ev(arg|disp|ent)|res)|ontinue|os)|d(el(com|ete|mod)|ialog|im|imtype|irinfo|irlist|ouble|up|upptr)|e(lse|nd|rr|x(ec|goto|ist|pf))|fo(nt|reach)|g(copy|et(ease|easef|key|path|str|time)|info|mode|osub|oto|radf|r(ect|oll|otate)|sel|square|zoom)|h(dc|instance|spstat|svcolor|spver|wnd)|i(f|nput|nstr|nt)|l(engt(h|h2|h3|h4)|i(bptr|mit|mitf|ne|stbox)|o(gf|gmes|op|oplev)|param|peek|poke)|m(call|ci|em(cpy|expand|file|set)|es|esbox|kdir|m(load|play|stop)|ous(e|ew|ex|ey)|ref)|new(com|lab|mod)|note(add|del|find|get|info|load|save|sel|unsel)|o(bj(enable|image|info|mode|prm|sel|size|skip)|n|n(click|cmd|error|exit|key))|p(al(color|ette)|eek|get|icload|oke|os|owf|rint|set)|querycom|r(andomize|e(draw|fdval|fstr|peat|turn)|nd|un)|s(arrayconv|creen|dim|endmsg|etease|in|ort(get|note|str|val)|plit|qrt|t(at|ick|r|rlen|rsize|op|r(f|mid|rep|trim))|ublev|ys(color|font|info))|tan|thismod|title|var(ptr|type|use)|w(ait|idth|inobj|param|peek|poke))(?=\(|\b)/gi, css: 'function' },
	      // preprocessor
	      { regex: /#\b(addition|aht|ahtmes|cfunc|cmd|cmpopt|comfunc|const|defcfunc|deffunc|define|else|endif|enum|epack|func|global|if|ifdef|ifndef|include|modcfunc|modfunc|modinit|modterm|module|pack|packopt|regcmd|runtime|undef|usecom|uselib)(?=\(|\b)/g, css: 'preprocessor' },
	      // macro
	      { regex: /\b(M_PI|_(_(date__|file__|hsp30__|hspver__|line__|time__)|break|continue|debug)|alloc|and|case|d(dim|efault|eg2rad|ir_(cmdline|cur|desktop|exe|mydoc|sys|tv|win)|o)|font_(antialias|bold|italic|normal|strikeout|underline)|for|ginfo_(accx|accy|accz|act|b|cx|cy|dispx|dispy|g|intid|mesx|mesy|mx|my|newid|paluse|r|sel|sizex|sizey|sx|sy|vx|vy|winx|winy|wx1|wx2|wy1|wy2)|gmode_(add|alpha|gdi|mem|pixela|rgb0|rgb0alpha|sub)|ldim|ms(gothic|mincho)|n(ext|ot|otemax|otesize)|obj(info_(bmscr|hwnd|mode)|mode_(guifont|normal|usefont))|or|rad2deg|screen_(fixedsize|frame|hide|normal|palette|tool)|swbreak|swend|switch|until|wend|while|xor)(?=\(|\b)/g, css: 'macro' },
	      // label
	      { regex: /^\s+\*[a-zA-Z].+?(\s|\n)|^\*[a-zA-Z].+?(\s|\n)/gm, css: 'label' }
	  ];
	  
	  this.forHtmlScript(SyntaxHighlighter.regexLib.aspScriptTags);
      };

      Brush.prototype	= new SyntaxHighlighter.Highlighter();
      Brush.aliases	= ['hsp', 'hsp3'];

      SyntaxHighlighter.brushes.Python = Brush;

      // CommonJS
      typeof(exports) != 'undefined' ? exports.Brush = Brush : null;
  })();
