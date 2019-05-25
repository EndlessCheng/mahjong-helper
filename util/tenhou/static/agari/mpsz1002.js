/////////////////////////////////////////////////////////////////////////////////////////////////////
// TENHOU.NET (C)C-EGG http://tenhou.net/
/////////////////////////////////////////////////////////////////////////////////////////////////////
var MPSZ={
	aka:true,
	fromHai136:function(hai136){
		var a=(hai136>>2);
		if (!this.aka) return ((a%9)+1) + "mpsz".substr(a/9,1);
		return (a<27 && (hai136%36)==16?"0":((a%9)+1)) + "mpsz".substr(a/9,1);
	},
	expand:function(t){
		return t
			.replace(/(\d)(\d{0,8})(\d{0,8})(\d{0,8})(\d{0,8})(\d{0,8})(\d{0,8})(\d{8})(m|p|s|z)/g,"$1$9$2$9$3$9$4$9$5$9$6$9$7$9$8$9")
			.replace(/(\d?)(\d?)(\d?)(\d?)(\d?)(\d?)(\d)(\d)(m|p|s|z)/g,"$1$9$2$9$3$9$4$9$5$9$6$9$7$9$8$9") // 57–‡˜A‘±‚ªãŒÀ 
			.replace(/(m|p|s|z)(m|p|s|z)+/g,"$1")
			.replace(/^[^\d]/,"");
	},
	contract:function(t){
		return t
			.replace(/\d(m|p|s|z)(\d\1)*/g,"$&:")
			.replace(/(m|p|s|z)([^:])/g,"$2")
			.replace(/:/g,"");
	},
	exsort:function(t){
		return t
			.replace(/(\d)(m|p|s|z)/g,"$2$1$1,")
			.replace(/00/g,"50")
			.split(",").sort().join("")
			.replace(/(m|p|s|z)\d(\d)/g,"$2$1");
	},
	exextract136:function(t){
		var s=t
			.replace(/(\d)m/g,"0$1")
			.replace(/(\d)p/g,"1$1")
			.replace(/(\d)s/g,"2$1")
			.replace(/(\d)z/g,"3$1");
		var i, c=new Array(136);
		for(i=0;i<s.length;i+=2){
			var n=s.substr(i,2), k=-1;
			if (n%10){
				var b=(9*Math.floor(n/10)+((n%10)-1))*4;
				k=(!c[b+3]?b+3:!c[b+2]?b+2:!c[b+1]?b+1:b);
			}else{
				k=(9*n/10+4)*4+0; // Ô
			}
			if (c[k]) document.write("err n="+n+" k="+k+"<br>");
			c[k]=1;
		}
		return c;
	},
	exextract34:function(t){
		var s=t
			.replace(/(\d)m/g,"0$1")
			.replace(/(\d)p/g,"1$1")
			.replace(/(\d)s/g,"2$1")
			.replace(/(\d)z/g,"3$1");
		var i, c=[0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0];
		for(i=0;i<s.length;i+=2){
			var n=s.substr(i,2), k=-1;
			if (n%10){
				k=9*Math.floor(n/10)+((n%10)-1);
			}else{
				k=9*n/10+4; // Ô
			}
			if (c[k]>4) document.write("err n="+n+" k="+k+"<br>");
			c[k]++;
		}
		return c;
	},
	compile136:function(c){
		var i, s="";
		for(i=0;i<136;++i) if (c[i]) s+=MPSZ.fromHai136(i);
		return s;
	}
};
