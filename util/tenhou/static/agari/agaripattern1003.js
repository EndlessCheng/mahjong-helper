/////////////////////////////////////////////////////////////////////////////////////////////////////
// TENHOU.NET (C)C-EGG http://tenhou.net/
/////////////////////////////////////////////////////////////////////////////////////////////////////
function AGARIPATTERN(){
	this.toitsu34=[-1,-1,-1,-1,-1,-1,-1];
	this.v=[{atama34:-1,mmmm35:0},{atama34:-1,mmmm35:0},{atama34:-1,mmmm35:0},{atama34:-1,mmmm35:0}]; // 一般形の面子の取り方は高々４つ 
	// mmmm35=( 21(順子)+34(暗刻)+34(槓子)+1(ForZeroInvalid) )*0x01010101 | 0x80808080(喰い)
}
AGARIPATTERN.prototype={
//	isKokushi:function(){return this.v[0].mmmm35==0xFFFFFFFF;},
//	isChiitoi:function(){return this.v[3].mmmm35==0xFFFFFFFF;},

	cc2m:function(c,d){
		return (c[d+0]<< 0)|(c[d+1]<< 3)|(c[d+2]<< 6)|
			(c[d+3]<< 9)|(c[d+4]<<12)|(c[d+5]<<15)|
			(c[d+6]<<18)|(c[d+7]<<21)|(c[d+8]<<24);
	},
	getAgariPattern:function(c,n){
		if (n!=34) return false;
		var e=this;
		var v=e.v;
		var j=(1<<c[27])|(1<<c[28])|(1<<c[29])|(1<<c[30])|(1<<c[31])|(1<<c[32])|(1<<c[33]);
		if (j>=0x10) return false; // 字牌が４枚 
		// 国士無双 // １４枚のみ
		if (((j&3)==2) && (c[0]*c[8]*c[9]*c[17]*c[18]*c[26]*c[27]*c[28]*c[29]*c[30]*c[31]*c[32]*c[33]==2)){
			var i,a=[0,8,9,17,18,26,27,28,29,30,31,32,33];
			for(i=0;i<13;++i) if (c[a[i]]==2) break;
			v[0].atama34=a[i];
			v[0].mmmm35=0xFFFFFFFF;
			return true;
		}
		if (j&2) return false; // 字牌が１枚 
		var ok=false;
		// 七対子 // １４枚のみ 
		if (!(j&10) && (
			(c[ 0]==2)+(c[ 1]==2)+(c[ 2]==2)+(c[ 3]==2)+(c[ 4]==2)+(c[ 5]==2)+(c[ 6]==2)+(c[ 7]==2)+(c[ 8]==2)+
			(c[ 9]==2)+(c[10]==2)+(c[11]==2)+(c[12]==2)+(c[13]==2)+(c[14]==2)+(c[15]==2)+(c[16]==2)+(c[17]==2)+
			(c[18]==2)+(c[19]==2)+(c[20]==2)+(c[21]==2)+(c[22]==2)+(c[23]==2)+(c[24]==2)+(c[25]==2)+(c[26]==2)+
			(c[27]==2)+(c[28]==2)+(c[29]==2)+(c[30]==2)+(c[31]==2)+(c[32]==2)+(c[33]==2))==7){
			v[3].mmmm35=0xFFFFFFFF;
			var i,n=0;
			for(i=0;i<34;++i) if (c[i]==2) e.toitsu34[n]=i, n+=1;
			ok=true;
			// 二盃口へ 
		}
		// 一般形
		var n00=c[ 0]+c[ 3]+c[ 6], n01=c[ 1]+c[ 4]+c[ 7], n02=c[ 2]+c[ 5]+c[ 8];
		var n10=c[ 9]+c[12]+c[15], n11=c[10]+c[13]+c[16], n12=c[11]+c[14]+c[17];
		var n20=c[18]+c[21]+c[24], n21=c[19]+c[22]+c[25], n22=c[20]+c[23]+c[26];
		var k0=(n00+n01+n02)%3;
		if (k0==1) return ok; // 余る
		var k1=(n10+n11+n12)%3;
		if (k1==1) return ok; // 余る
		var k2=(n20+n21+n22)%3;
		if (k2==1) return ok; // 余る
		if ((k0==2)+(k1==2)+(k2==2)+
			(c[27]==2)+(c[28]==2)+(c[29]==2)+(c[30]==2)+
			(c[31]==2)+(c[32]==2)+(c[33]==2)!=1) return ok; // 頭の場所は１つ 
		if (j&8){ // 字牌３枚
			if (c[27]==3) v[0].mmmm35<<=8, v[0].mmmm35|=21+27+1;
			if (c[28]==3) v[0].mmmm35<<=8, v[0].mmmm35|=21+28+1;
			if (c[29]==3) v[0].mmmm35<<=8, v[0].mmmm35|=21+29+1;
			if (c[30]==3) v[0].mmmm35<<=8, v[0].mmmm35|=21+30+1;
			if (c[31]==3) v[0].mmmm35<<=8, v[0].mmmm35|=21+31+1;
			if (c[32]==3) v[0].mmmm35<<=8, v[0].mmmm35|=21+32+1;
			if (c[33]==3) v[0].mmmm35<<=8, v[0].mmmm35|=21+33+1;
		}
		var n0=n00+n01+n02, kk0=(n00*1+n01*2)%3, m0=e.cc2m(c, 0);
		var n1=n10+n11+n12, kk1=(n10*1+n11*2)%3, m1=e.cc2m(c, 9);
		var n2=n20+n21+n22, kk2=(n20*1+n21*2)%3, m2=e.cc2m(c,18);
//		document.write("n="+n0+" "+n1+" "+n2+" k="+k0+" "+k1+" "+k2+" kk="+kk0+" "+kk1+" "+kk2+" mmmm="+v[0].mmmm35+"<br>");
		if (j&4){ // 字牌が頭 
			if (k0|kk0|k1|kk1|k2|kk2) return ok;
			if (c[27]==2) v[0].atama34=27;
			else if (c[28]==2) v[0].atama34=28;
			else if (c[29]==2) v[0].atama34=29;
			else if (c[30]==2) v[0].atama34=30;
			else if (c[31]==2) v[0].atama34=31;
			else if (c[32]==2) v[0].atama34=32;
			else if (c[33]==2) v[0].atama34=33;
			if (n0>=9){if (e.GetMentsu(1,m1) && e.GetMentsu(2,m2) && e.GetMentsu9Fin(0,m0)) return true;
			}else if (n1>=9){if (e.GetMentsu(2,m2) && e.GetMentsu(0,m0) && e.GetMentsu9Fin(1,m1)) return true;
			}else if (n2>=9){if (e.GetMentsu(0,m0) && e.GetMentsu(1,m1) && e.GetMentsu9Fin(2,m2)) return true;
			}else if (e.GetMentsu(0,m0) && e.GetMentsu(1,m1) && e.GetMentsu(2,m2)) return true; // 一意 
		}else if (k0==2){ // 萬子が頭 
			if (k1|kk1|k2|kk2) return ok;
			if (n0>=8){if (e.GetMentsu(1,m1) && e.GetMentsu(2,m2) && e.GetAtamaMentsu8Fin(kk0,0,m0)) return true;
			}else if (n1>=9){if (e.GetMentsu(2,m2) && e.GetAtamaMentsu(kk0,0,m0) && e.GetMentsu9Fin(1,m1)) return true;
			}else if (n2>=9){if (e.GetAtamaMentsu(kk0,0,m0) && e.GetMentsu(1,m1) && e.GetMentsu9Fin(2,m2)) return true;
			}else if (e.GetMentsu(1,m1) && e.GetMentsu(2,m2) && e.GetAtamaMentsu(kk0,0,m0)) return true; // 一意 
		}else if (k1==2){ // 筒子が頭 
			if (k2|kk2|k0|kk0) return ok;
			if (n1>=8){if (e.GetMentsu(2,m2) && e.GetMentsu(0,m0) && e.GetAtamaMentsu8Fin(kk1,1,m1)) return true;
			}else if (n2>=9){if (e.GetMentsu(0,m0) && e.GetAtamaMentsu(kk1,1,m1) && e.GetMentsu9Fin(2,m2)) return true;
			}else if (n0>=9){if (e.GetAtamaMentsu(kk1,1,m1) && e.GetMentsu(2,m2) && e.GetMentsu9Fin(0,m0)) return true;
			}else if (e.GetMentsu(2,m2) && e.GetMentsu(0,m0) && e.GetAtamaMentsu(kk1,1,m1)) return true; // 一意 
		}else if (k2==2){ // 索子が頭 
			if (k0|kk0|k1|kk1) return ok;
			if (n2>=8){if (e.GetMentsu(0,m0) && e.GetMentsu(1,m1) && e.GetAtamaMentsu8Fin(kk2,2,m2)) return true;
			}else if (n0>=9){if (e.GetMentsu(1,m1) && e.GetAtamaMentsu(kk2,2,m2) && e.GetMentsu9Fin(0,m0)) return true;
			}else if (n1>=9){if (e.GetAtamaMentsu(kk2,2,m2) && e.GetMentsu(0,m0) && e.GetMentsu9Fin(1,m1)) return true;
			}else if (e.GetMentsu(0,m0) && e.GetMentsu(1,m1) && e.GetAtamaMentsu(kk2,2,m2)) return true; // 一意 
		}
		v[0].mmmm35=0; // 一般形不発 
		return ok;
	},

// private:
	GetMentsu:function(col,m){ // ６枚以下は一意 
		var e=this;
		var mmmm=e.v[0].mmmm35;
		var i, a=(m&7), b=0, c=0;
		for(i=0;i<7;++i){
			switch(a){
			case 4:mmmm<<=16, mmmm|=((21+col*9+i+1)<<8) | (col*7+i+1), b+=1, c+=1;break;
			case 3:mmmm<<= 8, mmmm|=(21+col*9+i+1);break;
			case 2:mmmm<<=16, mmmm|=(col*7+i+1)*0x0101, b+=2, c+=2;break;
			case 1:mmmm<<= 8, mmmm|=(col*7+i+1), b+=1, c+=1;break;
			case 0:break;
			default:return false;
			}
			m>>=3, a=(m&7)-b, b=c, c=0;
		}
		if (a==3) mmmm<<=8, mmmm|=(21+col*9+7+1); else if (a) return false; // ⑧ 
		m>>=3, a=(m&7)-b;
		if (a==3) mmmm<<=8, mmmm|=(21+col*9+8+1); else if (a) return false; // ⑨ 
		e.v[0].mmmm35=mmmm;
//		DBGPRINT((_T("GetMentsu col=%d mmmm=%X\r\n"),col,mmmm));
		return true;
	},
	GetAtamaMentsu:function(nn,col,m){ // ５枚以下は一意 
		var e=this;
		var a=(7<<(24-nn*3));
		var b=(2<<(24-nn*3));
		if ((m&a)>=b && e.GetMentsu(col,m-b)) return e.v[0].atama34=col*9+8-nn, true;
		a>>=9, b>>=9;
		if ((m&a)>=b && e.GetMentsu(col,m-b)) return e.v[0].atama34=col*9+5-nn, true;
		a>>=9, b>>=9;
		if ((m&a)>=b && e.GetMentsu(col,m-b)) return e.v[0].atama34=col*9+2-nn, true;
		return false;
	},
	GetMentsu9:function(mmmm,col,m,v){ // const // ９枚以上
		// 面子選択は四連刻（１２枚）三連刻（９枚以上）しかない 
		var s=-1; // 三連刻
		var i, a=(m&7), b=0, c=0;
		for(i=0;i<7;++i){
			if (m==0x6DB) break; // 四連刻 // 三暗対々が高目 // １２枚のみ 
			switch(a){
			case 4:mmmm<<=8, mmmm|=(col*7+i+1), b+=1, c+=1; // nobreak // 平和二盃口が三暗刻より高目 
			case 3: // 帯幺九系が高目、ロン平和一盃口以外は三暗刻が高目 
				if (((m>>3)&7)>=3+b && ((m>>6)&7)>=3+c) s=i, b+=3, c+=3;// 三連刻 
				else mmmm<<=8, mmmm|=(21+col*9+i+1);
				break;
			case 2:mmmm<<=16, mmmm|=(col*7+i+1)*0x0101, b+=2, c+=2;break;
			case 1:mmmm<<= 8, mmmm|=(col*7+i+1), b+=1, c+=1;break;
			case 0:break;
			default:return 0;
			}
			m>>=3, a=(m&7)-b, b=c, c=0;
		}
		if (i<7){ // 四連刻を展開 
			v[0]=(21+col*9+i+1)*0x01010101 + 0x00010203;
			v[1]=(col*7+i+1+1)*0x010101 | (21+col*9+i+0+1)<<24;
			v[2]=(col*7+i+0+1)*0x010101 | (21+col*9+i+3+1)<<24;
			return 3;
		}
		if (a==3) mmmm<<=8, mmmm|=(21+col*9+7+1); else if (a) return 0; // ⑧ 
		m>>=3, a=(m&7)-b;
		if (a==3) mmmm<<=8, mmmm|=(21+col*9+8+1); else if (a) return 0; // ⑨ 

		if (s!=-1){ // 三連刻を展開 
			mmmm<<=24;
			v[0]=mmmm|((21+col*9+s+1)*0x010101 + 0x000102);
			v[1]=mmmm|((col*7+s+1)*0x010101);
			v[2]=0;
			return 2;
		}
		v[0]=mmmm, v[1]=v[2]=0;
		return 1;
	},
	GetMentsu9Fin:function(col,m){ // ９枚以上 
		var e=this;
		var v=e.v;
		var mm=[0,0,0];
		if (!e.GetMentsu9(v[0].mmmm35,col,m,mm)) return false;
		var n=0;
		if (mm[0]) v[n].atama34=v[0].atama34, v[n].mmmm35=mm[0], ++n;
		if (mm[1]) v[n].atama34=v[0].atama34, v[n].mmmm35=mm[1], ++n;
		if (mm[2]) v[n].atama34=v[0].atama34, v[n].mmmm35=mm[2], ++n;
//		document.write("GetMentsu9Fin col="+col+" n="+n+"<br>");
		return n!=0;
	},
	GetAtamaMentsu8Fin:function(nn,col,m){ // ８枚以上 
		var e=this;
		var v=e.v;
		var mmmm=v[0].mmmm35;
		var mm=[0,0,0];
		var a=(7<<(24-nn*3));
		var b=(2<<(24-nn*3));
		var n=0;
		if ((m&a)>=b && e.GetMentsu9(mmmm,col,m-b,mm)){
			if (mm[0]) v[n].atama34=col*9+8-nn, v[n].mmmm35=mm[0], ++n;
			if (mm[1]) v[n].atama34=col*9+8-nn, v[n].mmmm35=mm[1], ++n;
			if (mm[2]) v[n].atama34=col*9+8-nn, v[n].mmmm35=mm[2], ++n;
		}
		a>>=9, b>>=9;
		if ((m&a)>=b && e.GetMentsu9(mmmm,col,m-b,mm)){
			if (mm[0]) v[n].atama34=col*9+5-nn, v[n].mmmm35=mm[0], ++n;
			if (mm[1]) v[n].atama34=col*9+5-nn, v[n].mmmm35=mm[1], ++n;
			if (mm[2]) v[n].atama34=col*9+5-nn, v[n].mmmm35=mm[2], ++n;
		}
		a>>=9, b>>=9;
		if ((m&a)>=b && e.GetMentsu9(mmmm,col,m-b,mm)){
			if (mm[0]) v[n].atama34=col*9+2-nn, v[n].mmmm35=mm[0], ++n;
			if (mm[1]) v[n].atama34=col*9+2-nn, v[n].mmmm35=mm[1], ++n;
			if (mm[2]) v[n].atama34=col*9+2-nn, v[n].mmmm35=mm[2], ++n;
		}
//		document.write("GetAtamaMentsu8Fin col="+col+" n="+n+"<br>");
		return n!=0;
	}
};

/*
function getAgariPattern(c,n){
	if (n!=34) return;
	AGARIPATTERN.init();
	return AGARIPATTERN.getAgariPattern(c,n);
}
*/
