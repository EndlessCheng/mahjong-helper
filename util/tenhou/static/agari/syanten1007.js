/////////////////////////////////////////////////////////////////////////////////////////////////////
// TENHOU.NET (C)C-EGG http://tenhou.net/
/////////////////////////////////////////////////////////////////////////////////////////////////////
// REFERENCE
/////////////////////////////////////////////////////////////////////////////////////////////////////
// http://www.asamiryo.jp/fst13.html
/////////////////////////////////////////////////////////////////////////////////////////////////////


//function SYANTEN(a,n){}
//SYANTEN.prototype={
var SYANTEN={ // singleton
	n_eval:0,
	// input
	c:[0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0],
	// status
	n_mentsu :0,
	n_tatsu  :0,
	n_toitsu :0,
	n_jidahai:0, // １３枚にしてから少なくとも打牌しなければならない字牌の数 -> これより向聴数は下がらない 
	f_n4     :0, // 27bitを数牌、1bitを字牌で使用 
	f_koritsu:0, // 孤立牌
	// final result
	min_syanten:8,
	
	updateResult:function(){
		var e=this;
		var ret_syanten=8-e.n_mentsu*2-e.n_tatsu-e.n_toitsu;
		var n_mentsu_kouho=e.n_mentsu+e.n_tatsu;
		if (e.n_toitsu){
			n_mentsu_kouho+=e.n_toitsu-1;
		}else if (e.f_n4 && e.f_koritsu){
			if ((e.f_n4|e.f_koritsu)==e.f_n4) ++ret_syanten; // 対子を作成できる孤立牌が無い
		}
		if (n_mentsu_kouho>4) ret_syanten+=(n_mentsu_kouho-4);
		if (ret_syanten!=-1 && ret_syanten<e.n_jidahai) ret_syanten=e.n_jidahai;
		if (ret_syanten<e.min_syanten) e.min_syanten=ret_syanten;
	},


	// method
	init:function(a,n){
		var e=this;
		e.c=[0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0];
		// status
		e.n_mentsu =0;
		e.n_tatsu  =0;
		e.n_toitsu =0;
		e.n_jidahai=0;
		e.f_n4     =0;
		e.f_koritsu=0;
		// final result
		e.min_syanten=8;
		
		var c=this.c;
		if (n==136){
			for(n=0;n<136;++n) if (a[n]) ++c[n>>2];
		}else if (n==34){
			for(n=0;n<34;++n) c[n]=a[n];
		}else{
			for(n-=1;n>=0;--n) ++c[a[n]>>2];
		}
	},
	count34:function(){
		var c=this.c;
		return c[ 0]+c[ 1]+c[ 2]+c[ 3]+c[ 4]+c[ 5]+c[ 6]+c[ 7]+c[ 8]+
			c[ 9]+c[10]+c[11]+c[12]+c[13]+c[14]+c[15]+c[16]+c[17]+
			c[18]+c[19]+c[20]+c[21]+c[22]+c[23]+c[24]+c[25]+c[26]+
			c[27]+c[28]+c[29]+c[30]+c[31]+c[32]+c[33];
	},

	i_anko   :function(k){this.c[k]-=3, ++this.n_mentsu;},
	d_anko   :function(k){this.c[k]+=3, --this.n_mentsu;},
	i_toitsu :function(k){this.c[k]-=2, ++this.n_toitsu;},
	d_toitsu :function(k){this.c[k]+=2, --this.n_toitsu;},
	i_syuntsu:function(k){--this.c[k], --this.c[k+1], --this.c[k+2], ++this.n_mentsu;},
	d_syuntsu:function(k){++this.c[k], ++this.c[k+1], ++this.c[k+2], --this.n_mentsu;},
	i_tatsu_r:function(k){--this.c[k], --this.c[k+1], ++this.n_tatsu;},
	d_tatsu_r:function(k){++this.c[k], ++this.c[k+1], --this.n_tatsu;},
	i_tatsu_k:function(k){--this.c[k], --this.c[k+2], ++this.n_tatsu;},
	d_tatsu_k:function(k){++this.c[k], ++this.c[k+2], --this.n_tatsu;},
	i_koritsu:function(k){--this.c[k], this.f_koritsu|= (1<<k);},
	d_koritsu:function(k){++this.c[k], this.f_koritsu&=~(1<<k);},

	scanChiitoiKokushi:function(){
		var e=this;
		var syanten=e.min_syanten;
		var c=e.c;
		var n13= // 幺九牌の対子候補の数
			(c[ 0]>=2)+(c[ 8]>=2)+
			(c[ 9]>=2)+(c[17]>=2)+
			(c[18]>=2)+(c[26]>=2)+
			(c[27]>=2)+(c[28]>=2)+(c[29]>=2)+(c[30]>=2)+(c[31]>=2)+(c[32]>=2)+(c[33]>=2);
		var m13= // 幺九牌の種類数
			(c[ 0]!=0)+(c[ 8]!=0)+
			(c[ 9]!=0)+(c[17]!=0)+
			(c[18]!=0)+(c[26]!=0)+
			(c[27]!=0)+(c[28]!=0)+(c[29]!=0)+(c[30]!=0)+(c[31]!=0)+(c[32]!=0)+(c[33]!=0);
		var n7=n13+ // 対子候補の数
			(c[ 1]>=2)+(c[ 2]>=2)+(c[ 3]>=2)+(c[ 4]>=2)+(c[ 5]>=2)+(c[ 6]>=2)+(c[ 7]>=2)+
			(c[10]>=2)+(c[11]>=2)+(c[12]>=2)+(c[13]>=2)+(c[14]>=2)+(c[15]>=2)+(c[16]>=2)+
			(c[19]>=2)+(c[20]>=2)+(c[21]>=2)+(c[22]>=2)+(c[23]>=2)+(c[24]>=2)+(c[25]>=2);
		var m7=m13+ // 牌の種類数
			(c[ 1]!=0)+(c[ 2]!=0)+(c[ 3]!=0)+(c[ 4]!=0)+(c[ 5]!=0)+(c[ 6]!=0)+(c[ 7]!=0)+
			(c[10]!=0)+(c[11]!=0)+(c[12]!=0)+(c[13]!=0)+(c[14]!=0)+(c[15]!=0)+(c[16]!=0)+
			(c[19]!=0)+(c[20]!=0)+(c[21]!=0)+(c[22]!=0)+(c[23]!=0)+(c[24]!=0)+(c[25]!=0);
		{ // 七対子
			var ret_syanten=6-n7+(m7<7?7-m7:0);
			if (ret_syanten<syanten) syanten=ret_syanten;
		}
		{ // 国士無双
			var ret_syanten=13-m13-(n13?1:0);
			if (ret_syanten<syanten) syanten=ret_syanten;
		}
		return syanten;
	},
	removeJihai:function(nc){ // 字牌
		var e=this;
		var c=e.c;
		var j_n4=0; // 7bitを字牌で使用 
		var j_koritsu=0; // 孤立牌 
		var i;
		for(i=27;i<34;++i) switch(c[i]){
		case 4:++e.n_mentsu, j_n4|=(1<<(i-27)), j_koritsu|=(1<<(i-27)), ++e.n_jidahai;break;
		case 3:++e.n_mentsu;break;
		case 2:++e.n_toitsu;break;
		case 1:j_koritsu|=(1<<(i-27));break;
		}
		if (e.n_jidahai && (nc%3)==2) --e.n_jidahai;

		if (j_koritsu){ // 孤立牌が存在する
			e.f_koritsu|=(1<<27);
			if ((j_n4|j_koritsu)==j_n4) e.f_n4|=(1<<27); // 対子を作成できる孤立牌が無い 
		}
	},
	removeJihaiSanma19:function(nc){ // 字牌
		var e=this;
		var c=e.c;
		var j_n4=0; // 7+9bitを字牌で使用 
		var j_koritsu=0; // 孤立牌 
		var i;
		for(i=27;i<34;++i) switch(c[i]){
		case 4:++e.n_mentsu, j_n4|=(1<<(i-18)), j_koritsu|=(1<<(i-18)), ++e.n_jidahai;break;
		case 3:++e.n_mentsu;break;
		case 2:++e.n_toitsu;break;
		case 1:j_koritsu|=(1<<(i-18));break;
		}
		for(i=0;i<9;i+=8) switch(c[i]){
		case 4:++e.n_mentsu, j_n4|=(1<<i), j_koritsu|=(1<<i), ++e.n_jidahai;break;
		case 3:++e.n_mentsu;break;
		case 2:++e.n_toitsu;break;
		case 1:j_koritsu|=(1<<i);break;
		}
		if (e.n_jidahai && (nc%3)==2) --e.n_jidahai;

		if (j_koritsu){ // 孤立牌が存在する 
			e.f_koritsu|=(1<<27);
			if ((j_n4|j_koritsu)==j_n4) e.f_n4|=(1<<27); // 対子を作成できる孤立牌が無い 
		}
	},
	scanNormal:function(init_mentsu){
		var e=this;
		var c=e.c;
		e.f_n4|= // 孤立しても対子(雀頭)になれない数牌
			((c[ 0]==4)<< 0)|((c[ 1]==4)<< 1)|((c[ 2]==4)<< 2)|((c[ 3]==4)<< 3)|((c[ 4]==4)<< 4)|((c[ 5]==4)<< 5)|((c[ 6]==4)<< 6)|((c[ 7]==4)<< 7)|((c[ 8]==4)<< 8)|
			((c[ 9]==4)<< 9)|((c[10]==4)<<10)|((c[11]==4)<<11)|((c[12]==4)<<12)|((c[13]==4)<<13)|((c[14]==4)<<14)|((c[15]==4)<<15)|((c[16]==4)<<16)|((c[17]==4)<<17)|
			((c[18]==4)<<18)|((c[19]==4)<<19)|((c[20]==4)<<20)|((c[21]==4)<<21)|((c[22]==4)<<22)|((c[23]==4)<<23)|((c[24]==4)<<24)|((c[25]==4)<<25)|((c[26]==4)<<26);
		this.n_mentsu+=init_mentsu;
		e.Run(0);
	},

	Run:function(depth){ // ネストは高々１４回 
		var e=this;
		++e.n_eval;
		if (e.min_syanten==-1) return; // 和了は１つ見つければよい 
		var c=e.c;
		for(;depth<27&&!c[depth];++depth); // skip
		if (depth==27) return e.updateResult();

		var i=depth;
		if (i>8) i-=9;
		if (i>8) i-=9; // mod_9_in_27
		switch(c[depth]){
		case 4:
			// 暗刻＋順子|搭子|孤立
			e.i_anko(depth);
			if (i<7&&c[depth+2]){
				if (c[depth+1]) e.i_syuntsu(depth), e.Run(depth+1), e.d_syuntsu(depth); // 順子 
				e.i_tatsu_k(depth), e.Run(depth+1), e.d_tatsu_k(depth); // 嵌張搭子 
			}
			if (i<8&&c[depth+1]) e.i_tatsu_r(depth), e.Run(depth+1), e.d_tatsu_r(depth); // 両面搭子 
			e.i_koritsu(depth), e.Run(depth+1), e.d_koritsu(depth); // 孤立 
			e.d_anko(depth);
			// 対子＋順子系 // 孤立が出てるか？ // 対子＋対子は不可 
			e.i_toitsu(depth);
			if (i<7&&c[depth+2]){
				if (c[depth+1]) e.i_syuntsu(depth), e.Run(depth  ), e.d_syuntsu(depth); // 順子＋他 
				e.i_tatsu_k(depth), e.Run(depth+1), e.d_tatsu_k(depth); // 搭子は２つ以上取る必要は無い -> 対子２つでも同じ 
			}
			if (i<8&&c[depth+1]) e.i_tatsu_r(depth), e.Run(depth+1), e.d_tatsu_r(depth);
			e.d_toitsu(depth);
			break;
		case 3:
			// 暗刻のみ 
			e.i_anko(depth), e.Run(depth+1), e.d_anko(depth);
			// 対子＋順子|搭子 
			e.i_toitsu(depth);
			if (i<7&&c[depth+1]&&c[depth+2]){
				e.i_syuntsu(depth), e.Run(depth+1), e.d_syuntsu(depth); // 順子 
			}else{ // 順子が取れれば搭子はその上でよい 
				if (i<7&&c[depth+2]) e.i_tatsu_k(depth), e.Run(depth+1), e.d_tatsu_k(depth); // 嵌張搭子は２つ以上取る必要は無い -> 対子２つでも同じ 
				if (i<8&&c[depth+1]) e.i_tatsu_r(depth), e.Run(depth+1), e.d_tatsu_r(depth); // 両面搭子 
			}
			e.d_toitsu(depth);
			// 順子系 
			if (i<7&&c[depth+2]>=2&&c[depth+1]>=2) e.i_syuntsu(depth), e.i_syuntsu(depth), e.Run(depth), e.d_syuntsu(depth), e.d_syuntsu(depth); // 順子＋他 
			break;
		case 2:
			// 対子のみ 
			e.i_toitsu(depth), e.Run(depth+1), e.d_toitsu(depth);
			// 順子系
			if (i<7&&c[depth+2]&&c[depth+1]) e.i_syuntsu(depth), e.Run(depth), e.d_syuntsu(depth); // 順子＋他 
			break;
		case 1:
			// 孤立牌は２つ以上取る必要は無い -> 対子のほうが向聴数は下がる -> ３枚 -> 対子＋孤立は対子から取る 
			// 孤立牌は合計８枚以上取る必要は無い 
			if (i<6&&c[depth+1]==1&&c[depth+2]&&c[depth+3]!=4){ // 延べ単 
				e.i_syuntsu(depth), e.Run(depth+2), e.d_syuntsu(depth); // 順子＋他 
			}else{
//				if (n_koritsu<8) e.i_koritsu(depth), e.Run(depth+1), e.d_koritsu(depth);
				e.i_koritsu(depth), e.Run(depth+1), e.d_koritsu(depth);
				// 順子系
				if (i<7&&c[depth+2]){
					if (c[depth+1]) e.i_syuntsu(depth), e.Run(depth+1), e.d_syuntsu(depth); // 順子＋他 
					e.i_tatsu_k(depth), e.Run(depth+1), e.d_tatsu_k(depth); // 搭子は２つ以上取る必要は無い -> 対子２つでも同じ 
				}
				if (i<8&&c[depth+1]) e.i_tatsu_r(depth), e.Run(depth+1), e.d_tatsu_r(depth);
			}
			break;
		}
	}
};
function calcSyanten(a,n,bSkipChiitoiKokushi){
//	var e=new SYANTEN(a,n);
	var e=SYANTEN;
	e.init(a,n);
	var nc=e.count34();
	if (nc>14) return -2; // ネスト検査が爆発する 
	if (!bSkipChiitoiKokushi && nc>=13) e.min_syanten=e.scanChiitoiKokushi(nc); // １３枚より下の手牌は評価できない 
	e.removeJihai(nc);
//	e.removeJihaiSanma19(nc);
	var init_mentsu=Math.floor((14-nc)/3); // 副露面子を逆算 
	e.scanNormal(init_mentsu);
	return e.min_syanten;
}
function calcSyanten2(a,n){
//	var e=new SYANTEN(a,n);
	var e=SYANTEN;
	e.init(a,n);
	var nc=e.count34();
	if (nc>14) return undefined; // ネスト検査が爆発する 
	var syanten=[e.min_syanten,e.min_syanten];
	if (nc>=13) syanten[0]=e.scanChiitoiKokushi(nc); // １３枚より下の手牌は評価できない 
	e.removeJihai(nc);
//	e.removeJihaiSanma19(nc);
	var init_mentsu=Math.floor((14-nc)/3); // 副露面子を逆算 
	e.scanNormal(init_mentsu);
	syanten[1]=e.min_syanten;
	if (syanten[1]<syanten[0]) syanten[0]=syanten[1];
	return syanten;
}
