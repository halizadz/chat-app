# Analisis Bug dan Cacat Logika - Chat App

## 🔴 BUG KRITIS (Harus Diperbaiki)

### 1. **Backend: Auth Handler - Tidak Ada Pengecekan Duplikasi**
**Lokasi:** `backend/internal/handlers/auth.go:40-88`
**Masalah:** 
- Tidak ada pengecekan apakah email/username sudah terdaftar sebelum create
- Database akan return error generic jika duplicate, user tidak tahu field mana yang duplicate
- Error message tidak user-friendly

**Dampak:** User tidak tahu kenapa registrasi gagal (email atau username yang duplicate)

**Solusi:** Tambahkan pengecekan `FindByEmail` dan `FindByUsername` sebelum create, return error message yang jelas

---

### 2. **Backend: Chat Handler - User Bisa Menambahkan Dirinya Sendiri**
**Lokasi:** `backend/internal/handlers/chat.go:263-306`
**Masalah:**
- Tidak ada validasi untuk mencegah user menambahkan dirinya sendiri sebagai member
- User bisa menambahkan dirinya sendiri ke room

**Dampak:** Data inconsistency, user bisa menjadi member ganda

**Solusi:** Tambahkan pengecekan `if req.UserID == claims.UserID`

---

### 3. **Backend: Race Condition di FindOrCreatePrivateRoom**
**Lokasi:** `backend/internal/repository/room_repository.go:161-241`
**Masalah:**
- Jika 2 user create private room bersamaan, bisa terjadi duplicate rooms
- Tidak ada transaction lock atau unique constraint untuk private room

**Dampak:** Bisa terbuat 2 private room untuk 2 user yang sama

**Solusi:** Gunakan database-level lock atau unique constraint untuk private room

---

### 4. **Backend: SearchUsersByGmail - SQL Query Tidak Optimal**
**Lokasi:** `backend/internal/repository/user_repository.go:179-229`
**Masalah:**
- Query menggunakan 2x ILIKE: `WHERE email ILIKE $1 AND email ILIKE '%@gmail.com'`
- Bisa dioptimasi menjadi 1x ILIKE dengan pattern yang lebih baik

**Dampak:** Query lebih lambat dari seharusnya

**Solusi:** Optimasi query menjadi single ILIKE dengan pattern yang sudah include @gmail.com

---

### 5. **Backend: WebSocket - Tidak Ada Validasi Message Content**
**Lokasi:** `backend/internal/handlers/websocket.go:136-189`
**Masalah:**
- Tidak ada validasi panjang content message
- Tidak ada validasi untuk mencegah empty message
- Tidak ada sanitization untuk XSS

**Dampak:** 
- User bisa kirim message kosong atau sangat panjang
- Potensi XSS attack

**Solusi:** Tambahkan validasi content length dan sanitization

---

## 🟡 BUG MEDIUM (Perlu Diperbaiki)

### 6. **Frontend: useWebSocket - Reconnect Bisa Multiple Connections**
**Lokasi:** `frontend/src/hooks/useWebSocket.js:16-97`
**Masalah:**
- Reconnect logic tidak mengecek apakah sudah ada connection aktif
- Bisa terjadi multiple WebSocket connections untuk room yang sama
- Cleanup tidak sempurna

**Dampak:** Memory leak, duplicate messages

**Solusi:** Tambahkan pengecekan connection state sebelum reconnect

---

### 7. **Frontend: AddMemberModal - currentMembers Comparison Bug**
**Lokasi:** `frontend/src/components/modals/AddMemberModal.jsx:26-30`
**Masalah:**
- Comparison menggunakan `member.id === user.id || member.user_id === user.id`
- Tidak konsisten, bisa miss jika struktur data berbeda

**Dampak:** User yang sudah jadi member masih muncul di list

**Solusi:** Standardisasi comparison logic

---

### 8. **Frontend: chatStore - Tidak Ada Cleanup**
**Lokasi:** `frontend/src/store/chatStore.js`
**Masalah:**
- Tidak ada cleanup saat unmount
- WebSocket connections tidak di-close saat logout
- State tidak di-reset

**Dampak:** Memory leak, state tetap ada setelah logout

**Solusi:** Tambahkan cleanup function

---

## 🟢 BUG MINOR (Nice to Have)

### 9. **Backend: Auth - Tidak Ada Validasi Email Format**
**Lokasi:** `backend/internal/handlers/auth.go:40-88`
**Masalah:** Tidak ada validasi format email

**Dampak:** User bisa input email tidak valid

---

### 10. **Backend: Auth - Password Strength Tidak Divalidasi**
**Lokasi:** `backend/internal/handlers/auth.go:40-88`
**Masalah:** Tidak ada validasi minimal password length/strength

**Dampak:** Password lemah bisa digunakan

---

### 11. **Backend: Chat - Tidak Ada Validasi Room Type untuk AddMember**
**Lokasi:** `backend/internal/handlers/chat.go:263-306`
**Masalah:** User bisa add member ke private room (seharusnya hanya 2 member)

**Dampak:** Private room bisa punya lebih dari 2 member

---

### 12. **Frontend: NewChatModal - Gmail Detection Logic Tidak Sempurna**
**Lokasi:** `frontend/src/components/modals/NewChatModal.jsx:23-25`
**Masalah:** Hanya deteksi jika mengandung `@gmail.com`, tidak support username saja

**Dampak:** User harus ketik full email untuk search Gmail

---

## 📊 Ringkasan

- **Bug Kritis:** 5
- **Bug Medium:** 3  
- **Bug Minor:** 4
- **Total:** 12 bugs ditemukan

## 🎯 Prioritas Perbaikan

1. **PRIORITAS TINGGI:** Bug #1, #2, #3 (Data integrity & security)
2. **PRIORITAS SEDANG:** Bug #4, #5, #6, #7, #8 (Performance & UX)
3. **PRIORITAS RENDAH:** Bug #9-12 (Nice to have)
