/////////////////////////////////////////////////////////////////////////////////////////////////////
// TENHOU.NET (C)C-EGG http://tenhou.net/
/////////////////////////////////////////////////////////////////////////////////////////////////////
function $$(n,s){document.getElementById(n).innerHTML=s;}
/////////////////////////////////////////////////////////////////////////////////////////////////////
function MPSZ_exextract34(t){
	var n=parseInt(t.substr(0,1));
	return (n?n-1:4)+"mpsz".indexOf(t.substr(1,1))*9;
}
/////////////////////////////////////////////////////////////////////////////////////////////////////
function c_enum_machi34(c){
	var i, r=[];
	for(i=0;i<34;++i){
		if (c[i]>=4) continue;
		c[i]++; // 摸
		if (isAgari(c,34)) r.push(i);
		c[i]--;
	}
	return r;
}
function c136_to_c34(c136){
	var i, c=[0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0];
	for(i=0;i<136;++i) if (c136[i]) ++c[i>>2];
	return c;
}
/////////////////////////////////////////////////////////////////////////////////////////////////////
function img0(n,style,qt,q,da34,syanten){
	return "<a href=\"javascript:main('"+q+"', '"+qt+"');\" class=D onmouseover=\"daFocus(this,"+da34+","+syanten+");\" onmouseout=\"daUnfocus();\" >"+img1(n,style)+"</a>";
}
function img1(n,style){
	return "<img src=\"../static/img/agari/"+n+".gif\" border=0 "+(style?style:"")+" />";
}
function img2(n){
	return "<img src=\"../static/img/agari/"+n+".gif\" class=D />";
}
function img3(n,da34,qt,q){
	return "<a href=\"javascript:main('"+q+"', '"+qt+"');\" class=D onmouseover=\"daFocus(this,"+da34+");\" onmouseout=\"daUnfocus();\"><img src=\"../static/img/agari/"+n+".gif\" border=0 /></a>";
}
function sprintSyanten(n){
	return n==-1?"和了":n==0?"听牌":n+"向听";
}
function sprintSyanten2(a,b){
	if (b && a[0]!=a[1]) return "标准型"+sprintSyanten(a[0])+" / 一般型"+sprintSyanten(a[1]);
	return sprintSyanten(a[0]);
}

var lastE;
function daUnfocus(){
	if (lastE) lastE.style.backgroundColor="";
	$$("tips","");
}
function daFocus(a,da34,syanten){
	var e=document.getElementById("mda"+da34);
	if (e) e.style.backgroundColor="#CCCCCC";
	if (syanten!=undefined) $$("tips",sprintSyanten(syanten));
	lastE=e;
}
/////////////////////////////////////////////////////////////////////////////////////////////////////
function setup_yama(n){
	var i, yama=[];
	for(i=0;i<n;++i) yama.push(i);
	for(i=0;i<yama.length-1;++i){
		var j=i + Math.floor(Math.random()*(yama.length-i));
		var hai=yama[i];
		yama[i]=yama[j];
		yama[j]=hai;
	}
	return yama;
}
function sprintAgariPattern(c){
	var a2s=function(m){
		m=img1(MPSZ.fromHai136(m*4+1));
		return m+m;
	}
	var m2s=function(m){
		var kui=((m&0x80)!=0);
		m&=0x7F;
		if (m<21){
			m=Math.floor(m/7)*9 + (m%7);
			return img1(MPSZ.fromHai136(m*4+1))+img1(MPSZ.fromHai136(m*4+5))+img1(MPSZ.fromHai136(m*4+9));
		}else if (m<21+34){
			m-=21;
			m=img1(MPSZ.fromHai136(m*4+1));
			return m+m+m;
		}else if (m<21+34+34){
			m-=21+34;
			m=img1(MPSZ.fromHai136(m*4+1));
			return m+m+m+m;
		}
		return "";
	}

	var e=new AGARIPATTERN();
	if (!e.getAgariPattern(c,34)) return;
	var i, s="";
	for(i=0;i<4;++i){
		if (!e.v[i].mmmm35) continue;
		if (i==0 && e.v[0].mmmm35==0xFFFFFFFF){ // 国士無双
			s+="国士形和了 ";
			s+=a2s(e.v[i].atama34)+" ";
			var n,a=[0,8,9,17,18,26,27,28,29,30,31,32,33];
			for(n=0;n<13;++n) if (e.v[i].atama34!=a[n]) s+=img1(MPSZ.fromHai136(a[n]*4+1));
			s+="<br>";
		}else if (i==3 && e.v[3].mmmm35==0xFFFFFFFF){ // 七対子
			s+="七対形和了 ";
			s+=a2s(e.toitsu34[0])+" "+a2s(e.toitsu34[1])+" "+a2s(e.toitsu34[2])+" "+a2s(e.toitsu34[3])+" "+a2s(e.toitsu34[4])+" "+a2s(e.toitsu34[5])+" "+a2s(e.toitsu34[6]);
			s+="<br>";
		}else{
			var mentsu34=[
				((e.v[i].mmmm35>> 0)&0xFF)-1,
				((e.v[i].mmmm35>> 8)&0xFF)-1,
				((e.v[i].mmmm35>>16)&0xFF)-1,
				((e.v[i].mmmm35>>24)&0xFF)-1
			];
			s+="一般形和了 ";
			s+=a2s(e.v[i].atama34)+" "+m2s(mentsu34[3])+" "+m2s(mentsu34[2])+" "+m2s(mentsu34[1])+" "+m2s(mentsu34[0]);
			s+="<br>";
		}
	}
	return s;
}
/////////////////////////////////////////////////////////////////////////////////////////////////////
function printTehai(qt,q){
	var s="";
	s+="<hr size=1 color=#CCCCCC >";
	switch(qt.substr(0,1)){
	case "q":s+="标准型(包含七対国士)的计算结果"    +" / <a href=\""+"javascript:main(document.getElementById('tehaiInput').value, 'p');"+"\">标准型</a><br>";break;
	case "p":s+="一般型(不包含七対国士)的计算结果"+" / <a href=\""+"javascript:main(document.getElementById('tehaiInput').value, 'q');"+"\">一般型</a><br>";break;
	}
	var no14=(qt.substr(1,1)=="d");
	qt=qt.substr(0,1);

	q=MPSZ.expand(q);
	q=q.substr(0,14*2);
	var c=MPSZ.exextract136(q);
	var tsumo136=-1;
	while(tsumo136=Math.floor(Math.random()*136), c[tsumo136]);
	var qn3=(Math.floor(q.length/2)%3);
	if (qn3!=2 && !no14){ // ツモ
		c[tsumo136]=1;
		q+=MPSZ.fromHai136(tsumo136);
	}
	var restc=function(v,c){ // 受け入れ枚数を数える
		var i,n=0;
		for(i=0;i<v.length;++i) n+=4-c[v[i]];
		return n;
	}

	c=c136_to_c34(c);
	var tehai="";
	var syanten_org01=calcSyanten2(c,34);
	tehai+=sprintSyanten2(syanten_org01,q.length==14*2);
	tehai+="("+Math.floor(q.length/2)+"枚)";
	if (syanten_org01[0]==-1) tehai+=" / <a href=\"javascript:main('', _qt);\" >生成新手牌</a>";
	tehai+="<br/>";

	var syanten_org=(qt=="q"?syanten_org01[0]:syanten_org01[1]);
	var i,j;
	var v=new Array(34+1);
	if (syanten_org==0 && qn3==1 && no14){
		i=34;
		v[i]=c_enum_machi34(c);
		if (v[i].length) v[i]={da:i,n:restc(v[i],c),v:v[i]};
	}else if (syanten_org<=0){
		for(i=0;i<34;++i){
			if (!c[i]) continue;
			c[i]--; // 打
			v[i]=c_enum_machi34(c);
			c[i]++;
			if (v[i].length) v[i]={da:i,n:restc(v[i],c),v:v[i]};
		}
	}else if (qn3==2 || (qn3==1 && !no14)){
		var tt=(new Date()).getTime();
		for(i=0;i<34;++i){
			if (!c[i]) continue;
			c[i]--; // 打
			v[i]=[];
			for(j=0;j<34;++j){
				if (i==j || c[j]>=4) continue;
				c[j]++; // 摸
				if (calcSyanten(c,34,qt=="p")==syanten_org-1) v[i].push(j);
				c[j]--;
			}
			c[i]++;
			if (v[i].length) v[i]={da:i,n:restc(v[i],c),v:v[i]};
		}
		tt=(new Date()).getTime()-tt;
//		s+="<br>time="+tt+" eval="+SYANTEN.n_eval;
	}else{
		i=34;
		v[i]=[];
		for(j=0;j<34;++j){
			if (c[j]>=4) continue;
			c[j]++; // 摸
			if (calcSyanten(c,34,qt=="p")==syanten_org-1) v[i].push(j);
			c[j]--;
		}
		if (v[i].length) v[i]={da:i,n:restc(v[i],c),v:v[i]};
	}
	
	var trtd=[];
	for(i=0;i<q.length;i+=2){
		var da=q.substr(i,2);
		var da34=MPSZ_exextract34(da);
		var t=MPSZ.contract(MPSZ.exsort(q.replace(da,"")));
		var syanten=syanten_org+1;
		var vd=v[da34];
		if (vd && vd.n){
			syanten=(syanten_org==-1?0:syanten_org);
			if (vd.q==undefined) trtd.push(vd);
			vd.q=t; // 上書きで赤打の優先度が下がる
		}
		if (qn3==2) t+=MPSZ.fromHai136(tsumo136);
		tehai+=(qn3==2 || (qn3!=2 && !no14)?img0:img1)(da,((i%3)==2 && i==q.length-2?" hspace=3 ":""),qt,t,da34,syanten);
	}
	if (v[34] && v[34].n){
		v[34].q=MPSZ.contract(q);
		trtd.push(v[34]);
		tehai+="<br><br><a href=\"?"+qt+"="+v[34].q+"\">随机自摸</a>";
	}

	trtd.sort(function(a,b){return b.n-a.n;});
	s+="<table cellpadding=2 cellspacing=0 >";
	var tdisp=(syanten_org<=0?"待":"摸");
	for(i=0;i<trtd.length;++i){
		var da34=trtd[i].da;
		s+="<tr id=mda"+da34+" ><td>";
		if (da34<34) s+="打</td><td>"+img2(MPSZ.fromHai136(da34*4+1))+"</td><td>";
		s+=tdisp+"[</td><td>";
		var j, v=trtd[i].v, qq=trtd[i].q;
		for(j=0;j<v.length;++j){
			var tsumo=MPSZ.fromHai136(v[j]*4+1);
			s+=img3(tsumo,da34,qt,qq+tsumo);
		}
		s+="</td><td>"+trtd[i].n+"枚</td><td>]</td></tr>";
	}
	s+="</table>";
	if (syanten_org01[0]==-1){
		s+="<hr size=1 color=#CCCCCC >";
		s+=sprintAgariPattern(c);
	}
	$$("tehai",tehai);
	$$("tips","");
	$$("m2",s);
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

function main(q, qt){
	document.getElementById("tehaiInput").value=q;
	if (!q.length){
		var i, s="", c=Math.floor(Math.random()*3);
		s+="<table cellpadding=0 cellspacing=0 ><tr><td>";
		for(i=9;i>=0;--i){
			var yama=setup_yama(i>3?4*34:4*9);
			var tehai=yama.splice(0,i==9?4:i==8?7:i==7?10:13).sort(function(a,b){return a-b;});
			var tsumo136=yama.splice(0,1)[0];
			var j,q="";
			for(j=0;j<tehai.length;++j) q+=MPSZ.fromHai136(tehai[j]+(i>3?0:4*9*c));
			if (tsumo136!=-1) q+=MPSZ.fromHai136(tsumo136+(i>3?0:4*9*c)), tehai.push(tsumo136);
			var syanten01=calcSyanten2(tehai,tehai.length);
			s+="■<a href=\"javascript:main('"+MPSZ.contract(q)+"', _qt);\" class=D >"
			s+=sprintSyanten2(syanten01,tehai.length==14);
			s+="<br>";
			for(j=0;j<q.length;j+=2) s+=img1(q.substr(j,2),(j%6)==2&&j==q.length-2?" hspace=3 ":"");
			s+="</a>";
			s+="<br><br>";
			if (i==4) document.getElementById("tehaiInput").value=MPSZ.contract(q);
		}
		s+="</td></tr></table>";
		$$("tehai",s);
	}else if (q.match(/^(\d+m|\d+p|\d+s|[1234567]+z)*$/)==null){
		$$("tehai","<font color=#FF0000 >INVALID QUERY</font>");
	}else if (q.length){
		printTehai(qt,q);
	}
}

var _qt="q";
var _q=getQueryStringByName("q");
if(!_q) {
	_q=getQueryStringByName("p");
	_qt="p";
}
main(_q, _qt);
