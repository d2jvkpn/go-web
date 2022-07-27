### Constraints
---

#### 1. Response JSON Fields
- code: int64, must;
- msg: string, must & default "" or "ok";
- data: object(map), must & default {};
  - respone items:
    - total: int64;
    - items: []any, must & default [] & not nil(null);
  - key-value pairs

#### 2. Context Keys
- : abcDef
- : AbcDef
- : _AbcDef

#### 3. Status Codes(HTTP & response.Code)


#### 4. Logger

#### 4.2 required fields
- traceId:
- spanId:
- eventId(requestId):

#### 4.2 keys
- apis: api/web, api/mobile...;
- internals: internal/cron, internal/rpc, internal/panic, internal/conflict...;
- business: biz/xxxx;
- third parties: thirdParty/CloudProvider::Name...;
