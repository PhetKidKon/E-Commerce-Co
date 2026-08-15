# E-Commerce (Microservices)

Skeleton project. Backend เป็น Go microservices หลังประตู API gateway
(**nginx → Kong db-less → services**) และ Frontend เป็น React + TypeScript

```
E-Commence[Co]/
├── Backend/
│   ├── go.work                  # Go workspace (Common + 4 services)
│   ├── docker-compose.yml       # รัน 4 services บน network "ecommerce" (โมเดล B)
│   ├── Common-Service/          # shared library (response helper) — ไม่มี main(), import อย่างเดียว
│   ├── User-Service/            # :8081  + User Service
│   ├── Product-Service/         # :8082  + Product Service
│   ├── Cart-Service/            # :8083  + Cart Service
│   ├── Inventory-Service/       # :8084  + Inventory Service
│   ├── Gateway/                 # API gateway = nginx (edge) + Kong (db-less)
│   │   ├── docker-compose.yml   #   รัน nginx + kong ด้วยกัน
│   │   ├── Dockerfile           #   Kong image
│   │   ├── kong/kong.yml        #   route /api/* -> แต่ละ service
│   │   ├── nginx/               #   nginx.conf: edge reverse proxy -> kong
│   │   ├── custom-auth/         #   ทำทีหลัง (Lua?)
│   │   └── deployment/          #   ทำทีหลัง (Jenkin?)
└── Frontend/                    # React + TypeScript (Vite) :5173
```

## Ports

| Component                    | Port      | หมายเหตุ                        |
|-----------------------------|-----------|--------------------------------|
| Frontend (Vite)             | 5173      |                                |
| User/Product/Cart/Inventory | 8081-8084 | เรียกตรงได้ตอนรัน local (โมเดล A) |
| nginx (edge)                | 8080      | จุด entry เดียว (โมเดล B)        |
| Kong proxy / admin          | 8000/8001 | ภายใน / dev เท่านั้น (โมเดล B)   |

---


## รันแบบ Local

รัน service ตรงๆ บนเครื่องด้วย `go run` แก้โค้ดแล้วรันใหม่ได้ทันที ไม่ต้อง build image
เหมาะกับพัฒนา/debug เพราะ log เด้งขึ้น terminal และต่อ debugger ได้ตรงๆ

```bash
# Backend — เปิดคนละ terminal ต่อ service
cd Backend
go run ./User-Service        # :8081
go run ./Product-Service     # :8082
go run ./Cart-Service        # :8083
go run ./Inventory-Service   # :8084

# Frontend
cd Frontend
npm install                  # ครั้งแรกครั้งเดียว
npm run dev                  # http://localhost:5173
```

ทดสอบ (เรียก service ตรงๆ ไม่ผ่าน gateway):

```bash
curl http://localhost:8081/health
curl http://localhost:8082/health
```

> โมเดลนี้ยังไม่ใช้ gateway — ยิงเข้า service ที่ port ของมันตรงๆ
> ถ้าต้องการทดสอบผ่าน gateway (route รวม /api/*) ให้ใช้โมเดล B

---

## รันโดยใช้ Docker

รันทุกอย่างเป็น container เหมือน prod — เข้าทาง gateway ทางเดียว (`:8080`)
เหมาะกับ integration test / demo / เช็คว่ารันใน container ได้จริง ต้องเปิด Docker ก่อน

```bash
# 1) สร้าง shared network (ทำครั้งเดียว)
#    จำเป็นเพราะ compose แยกเป็น 2 ไฟล์ (Backend + Gateway) ต้องใช้ network ร่วมกัน
docker network create ecommerce

# 2) รัน backend services
cd Backend
docker compose up --build -d

# 3) รัน gateway (nginx + kong)   *ต้องรันหลัง services*
cd Gateway
docker compose up --build -d
```

ทดสอบ (ผ่าน gateway):

```bash
curl http://localhost:8080/health         # nginx edge
curl http://localhost:8080/api/users      # nginx -> kong -> user-service
```

หยุด:

```bash
cd Backend/Gateway && docker compose down
cd Backend         && docker compose down
```

---

## หมายเหตุ

- ตอนนี้ service ทุกตัวเป็น **skeleton** มีแค่ `/health` → เรียก `/api/*` จะได้ **404**
  จนกว่าจะใส่ handler (route ทำงานถึง service แล้ว แต่ปลายทางยังไม่มีของให้ตอบ)
- โครง layer (`internal/{model,repository,service,handler}`) วางไว้แล้วแต่ยังว่าง
- โค้ด service ชุดเดียวใช้ได้ทั้งโมเดล A และ B — ต่างแค่ "วิธีรัน" ไม่ต้องแก้โค้ด
- Kong เป็น db-less — แก้ routing ที่ `Gateway/kong/kong.yml` ได้ตรงๆ (แก้แล้ว restart container พอ ไม่ต้อง rebuild)
- เปิด custom Lua plugin ทีหลัง: ปลด comment `KONG_PLUGINS: "bundled,custom-auth"` ใน `Gateway/docker-compose.yml`
- module path: `github.com/kidkon/ecommerce/*` (resolve เป็น local ผ่าน go.work ไม่ได้โหลดจาก GitHub)
