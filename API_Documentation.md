# **NavMate API Specification**

เอกสารนี้สรุปรายละเอียดของ API Endpoints สำหรับแอปพลิเคชัน NavMate

**Base URL:** `http://localhost:8080`

**Authentication:** การเรียก API ส่วนใหญ่ในกลุ่ม `/v1` จำเป็นต้องมีการยืนยันตัวตนผ่าน JWT โดยส่ง `Authorization` Header ในรูปแบบ `Bearer <YOUR_JWT_TOKEN>`

-----

## **1. Health Check**

Endpoints สำหรับตรวจสอบสถานะของเซิร์ฟเวอร์

### **GET /health**

  * **Description:** ตรวจสอบว่าเซิร์ฟเวอร์ทำงานปกติหรือไม่
  * **Authentication:** ไม่จำเป็น
  * **Success Response (200 OK):**
    ```json
    {
      "status": "ok",
      "timestamp": "2025-09-05T04:16:00Z"
    }
    ```

-----

## **2. Authentication**

Endpoints สำหรับการจัดการผู้ใช้และการยืนยันตัวตน

### **POST /v1/auth/signup**

  * **Description:** สมัครสมาชิกใหม่ด้วยอีเมลและรหัสผ่าน
  * **Authentication:** ไม่จำเป็น
  * **Request Body:**
    ```json
    {
      "email": "test@example.com",
      "password": "password123"
    }
    ```
  * **Success Response (201 Created):**
    ```json
    {
      "id": 1,
      "email": "test@example.com"
    }
    ```

### **POST /v1/auth/login**

  * **Description:** เข้าสู่ระบบเพื่อรับ JWT Token
  * **Authentication:** ไม่จำเป็น
  * **Request Body:**
    ```json
    {
      "email": "test@example.com",
      "password": "password123"
    }
    ```
  * **Success Response (200 OK):**
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```

### **GET /v1/me**

  * **Description:** ดึงข้อมูลโปรไฟล์ของผู้ใช้ที่กำลังล็อกอินอยู่
  * **Authentication:** **จำเป็น**
  * **Success Response (200 OK):**
    ```json
    {
      "user_id": 1,
      "email": "test@example.com"
    }
    ```

### **GET /auth/google/login**

  * **Description:** เริ่มกระบวนการล็อกอินด้วย Google โดยจะ Redirect ไปยังหน้า Google Authentication
  * **Authentication:** ไม่จำเป็น

### **GET /auth/google/callback**

  * **Description:** Endpoint ที่ Google จะเรียกกลับมาหลังจากการยืนยันตัวตนสำเร็จ ระบบจะสร้างผู้ใช้ (หากยังไม่มี) และคืน JWT Token กลับไป
  * **Authentication:** ไม่จำเป็น
  * **Success Response (200 OK):**
    ```json
    {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "email": "user.from.google@gmail.com",
        "provider": "google"
    }
    ```

-----

## **3. Trip Planning**

Endpoints สำหรับการวางแผนการเดินทาง

### **POST /v1/trips/plan**

  * **Description:** สร้างแผนการเดินทางใหม่โดยระบุต้นทางและปลายทาง ระบบจะคืนตัวเลือกการเดินทาง (Itineraries) ที่เป็นไปได้กลับมา
  * **Authentication:** **จำเป็น**
  * **Request Body:**
    ```json
    {
      "origin": "Siam Paragon",
      "destination": "Central World",
      "depart_at": "2025-09-05T10:00:00Z"
    }
    ```
  * **Success Response (200 OK):**
    ```json
    {
      "plan_id": 1,
      "options": [
        { "itinerary_id": 1, "mode_mix": "WALK+TRANSIT", "total_minutes": 42, "rough_cost_cents": 3000 },
        { "itinerary_id": 2, "mode_mix": "RIDE", "total_minutes": 18, "rough_cost_cents": 12000 },
        { "itinerary_id": 3, "mode_mix": "WALK+TRANSIT+RIDE", "total_minutes": 28, "rough_cost_cents": 9000 }
      ]
    }
    ```

### **GET /v1/trips/plans/:id**

  * **Description:** ดึงข้อมูลแผนการเดินทางตาม `plan_id`
  * **Authentication:** **จำเป็น**
  * **Success Response (200 OK):**
    ```json
    {
        "id": 1,
        "origin": "Siam Paragon",
        "destination": "Central World",
        "status": "planned",
        "selected_itinerary_id": null,
        "itinerary_count": 3
    }
    ```

### **POST /v1/trips/plans/:id/select**

  * **Description:** เลือกตัวเลือกการเดินทาง (Itinerary) สำหรับแผนที่กำหนด
  * **Authentication:** **จำเป็น**
  * **Request Body:**
    ```json
    {
      "itinerary_id": 2
    }
    ```
  * **Success Response:** `204 No Content`

-----

## **4. Booking**

Endpoints สำหรับการจองการเดินทาง (เช่น เรียกรถ)

### **POST /v1/bookings**

  * **Description:** สร้างการจองสำหรับแผนการเดินทางที่เลือกไว้
  * **Authentication:** **จำเป็น**
  * **Request Body:**
    ```json
    {
      "plan_id": 1
    }
    ```
  * **Success Response (200 OK):**
    ```json
    {
      "booking_id": 1,
      "status": "confirmed",
      "eta_minutes": 10,
      "fare_cents": 12100
    }
    ```

### **GET /v1/bookings/:id**

  * **Description:** ดึงข้อมูลการจองตาม `booking_id`
  * **Authentication:** **จำเป็น**
  * **Success Response (200 OK):**
    ```json
    {
      "ID": 1,
      "plan_id": 1,
      "itinerary_id": 2,
      "provider": "RideNow",
      "status": "confirmed",
      "eta_minutes": 10,
      "fare_cents": 12100,
      "payment_id": null,
      "quote_id": "",
      "external_ref": "",
      "created_at": "2025-09-05T04:16:00Z",
      "updated_at": "2025-09-05T04:16:00Z"
    }
    ```

-----

## **5. Payment**

Endpoints สำหรับการจัดการการชำระเงิน

### **POST /v1/payments/authorize**

  * **Description:** ทำการกันวงเงิน (Authorize) สำหรับการจอง
  * **Authentication:** **จำเป็น**
  * **Request Body:**
    ```json
    {
      "booking_id": 1,
      "amount_cents": 12100
    }
    ```
  * **Success Response (200 OK):**
    ```json
    {
      "payment_id": 1,
      "status": "authorized",
      "external_ref": "stub-auth-123"
    }
    ```

### **POST /v1/payments/:id/capture**

  * **Description:** ยืนยันการตัดเงิน (Capture) ที่ได้ทำการกันวงเงินไว้
  * **Authentication:** **จำเป็น**
  * **Success Response (200 OK):**
    ```json
    {
        "payment_id": 1,
        "status": "captured"
    }
    ```

### **POST /v1/payments/:id/refund**

  * **Description:** ขอคืนเงิน (Refund) สำหรับการชำระเงินที่สำเร็จแล้ว
  * **Authentication:** **จำเป็น**
  * **Success Response (200 OK):**
    ```json
    {
        "payment_id": 1,
        "status": "refunded"
    }
    ```

-----

## **6. Safety**

Endpoints สำหรับฟีเจอร์ติดตามความปลอดภัย

### **POST /v1/safety/session**

  * **Description:** เริ่ม Safety Session สำหรับแผนการเดินทาง
  * **Authentication:** **จำเป็น**
  * **Request Body:**
    ```json
    {
      "plan_id": 1,
      "interval_min": 15
    }
    ```
  * **Success Response (200 OK):**
    ```json
    {
      "session_id": 1,
      "share_url": "/safety/s/RANDOM_TOKEN_STRING",
      "next_due": "2025-09-05T04:45:00Z",
      "interval_minutes": 15
    }
    ```

### **POST /v1/safety/heartbeat/ack**

  * **Description:** ยืนยันความปลอดภัย (Check-in) สำหรับ Safety Session ที่ทำงานอยู่
  * **Authentication:** **จำเป็น**
  * **Request Body:**
    ```json
    {
      "session_id": 1
    }
    ```
  * **Success Response (200 OK):**
    ```json
    {
      "next_due": "2025-09-05T05:00:00Z",
      "status": "acknowledged"
    }
    ```

### **POST /v1/safety/sos**

  * **Description:** ส่งสัญญาณขอความช่วยเหลือฉุกเฉิน (SOS)
  * **Authentication:** **จำเป็น**
  * **Request Body:**
    ```json
    {
      "plan_id": 1,
      "location": "Lat: 13.746, Lon: 100.535",
      "message": "I need help!"
    }
    ```
  * **Success Response (200 OK):**
    ```json
    {
      "status": "SOS triggered",
      "plan_id": 1,
      "timestamp": "2025-09-05T04:16:00Z",
      "message": "Emergency services and contacts will be notified"
    }
    ```

### **GET /safety/s/:token**

  * **Description:** (Public) Endpoint สำหรับดูสถานะของ Safety Session ผ่าน Share Token
  * **Authentication:** ไม่จำเป็น
  * **Success Response (200 OK):**
    ```json
    {
      "plan_id": 1,
      "origin": "Siam Paragon",
      "destination": "Central World",
      "started_at": "2025-09-05T04:30:00Z",
      "next_due": "2025-09-05T04:45:00Z",
      "active": true,
      "interval_minutes": 15,
      "last_heartbeat": {
          "due_at": "2025-09-05T04:45:00Z",
          "status": "due",
          "acked_at": null
      }
    }
    ```