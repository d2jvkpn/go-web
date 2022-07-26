### Api constraints
---

#### 1. Response JSON Fields
- code: int64, must;
- msg: string, must & default "" or "ok";
- data: object(map), must & default {};
  - respone items:
    - total: int64;
    - items: []any, must & default [] & not nil(null);
  - key-value pairs

#### 2. HTTP Status Codes
