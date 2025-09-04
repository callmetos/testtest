# ---- Build Stage ----
FROM golang:1.23-bookworm AS build
WORKDIR /app

# ติดตั้ง git และ ca-certificates ที่จำเป็นสำหรับการดึง Go modules
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy ไฟล์ go.mod และ go.sum เข้ามาก่อน เพื่อใช้ประโยชน์จาก Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy โค้ดของโปรเจกต์ทั้งหมดเข้ามา
COPY . .

# Build โปรแกรม Go ให้เป็น static binary สำหรับ Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/main.go

# ---- Run Stage ----
# ใช้ 'distroless' image ซึ่งมีขนาดเล็กและปลอดภัยสูง
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /home/nonroot

# Copy ไฟล์ binary ที่ build เสร็จแล้วจาก stage ก่อนหน้า
COPY --from=build /out/api ./api

# ระบุว่า container จะทำงานที่พอร์ต 8080
EXPOSE 8080

# สั่งให้ container รันด้วย user ที่ไม่มีสิทธิ์ root เพื่อความปลอดภัย
USER nonroot:nonroot

# คำสั่งที่จะรันเมื่อ container เริ่มทำงาน
ENTRYPOINT ["./api"]