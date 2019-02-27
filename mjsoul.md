# 雀魂协议
  
## 协议基本

    雀魂协议 使用 protobuff 
    具体协议 https://majsoul.union-game.com/0/v0.4.243.w/res/proto/liqi.json

## 第一位

    raw[0] === 1 表示RESPONSE,REQUEST

    raw[0] === 3 表示NOTIFY

## NOTIFY

### 去掉第一位

    raw = raw.slice(1)

### 处理协议

#### 外层协议

      "Wrapper":{
          "fields":{
              "name":{
                  "type":"string",
                  "id":1
              },
              "data":{
                  "type":"bytes",
                  "id":2
              }
          }
      },

     获得 方法名字 name 和 数据 data 

##### 常用协议

      NotifyPlayerLoadGameReady
      ActionPrototype
      NotifyGameEndResult

#### 解析 ActionPrototype

##### 协议

      "ActionPrototype":{
          "fields":{
              "step":{
                  "type":"uint32",
                  "id":1
              },
              "name":{
                  "type":"string",
                  "id":2
              },
              "data":{
                  "type":"bytes",
                  "id":3
              }
          }
      },

#### 常用协议
    
    ActionNewRound
    ActionDiscardTile
    ActionDealTile
    ActionChiPengGang
    ActionAnGangAddGang

##### ActionNewRound
  
    {"chang":0,"ju":0,"ben":0,"tiles":["3m","6m","6m","7m","7m","8m","9m","3p","8p","1s","2s","4s","9s","1z"],"dora":"4s","scores":[25000,25000,25000,25000],"operation":{"seat":0,"operation_list":[{"type":1}],"time_add":20000,"time_fixed":8000},"liqibang":0,"al":false,"md5":"005451C30A3B549D9B7AF9ED09CDFED5","left_tile_count":69}

##### ActionDealTile
  
###### 其他人

      {"seat":1,"left_tile_count":68,"zhenting":false}

###### 自己

      {"seat":1,"left_tile_count":68,"zhenting":false}

##### ActionDiscardTile
  
      {"seat":0,"tile":"1z","is_liqi":false,"moqie":true,"zhenting":false,"is_wliqi":false}




