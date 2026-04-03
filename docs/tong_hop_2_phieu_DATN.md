# Tong hop 2 phieu DATN theo tien do project

Tai lieu nay duoc tong hop tu:
- `/Users/hant/golang/doan/PROJECT_TASKS.md`
- `/Users/hant/golang/doan/README.md`
- `/Users/hant/golang/doan/docs/bao cao tuan 3-4.docx` (noi dung text)
- lich su commit `git log`

Luu y:
- Cac truong hanh chinh khong xac dinh chac chan tu repo duoc de dang `[Dien ...]`.
- Thong tin sinh vien suy ra tu `README.md`: `Nguyen The Ha - 61165 - CS2`.
- Noi dung duoi day duoc viet de ban chep vao dung 2 bieu mau DATN-01 va DATN-02.

---

## PHIEU 1 - MAU DATN-01

**NHIEM VU DO AN TOT NGHIEP**

| Muc | Noi dung de dien |
|---|---|
| Ho va ten sinh vien | Nguyen The Ha |
| Ma SV | 61165 |
| Lop | CS2 |
| Khoa | [Dien khoa] |
| Nganh/Chuyen nganh | [Dien nganh/chuyen nganh] |
| He dao tao | [Dien he dao tao] |
| Ten de tai DATN | Xay dung he thong quan ly trung tam day them tich hop AI ho tro kiem soat chat luong dao tao |
| Khoa | [Dien khoa quan ly de tai] |
| Nhom chuyen mon | [Dien nhom chuyen mon] |
| GV huong dan | [Dien ten GVHD] |
| Email/DT | [Dien email/so dien thoai GVHD] |
| GV dong huong dan (neu co) | [De trong neu khong co] |
| Email/DT | [Dien neu co] |
| Doanh nghiep/don vi phoi hop (neu co) | [De trong neu khong co] |
| Dia diem thuc hien | /Users/hant/golang/doan va moi truong thuc nghiem local |
| Thoi gian thuc hien | 15 tuan (26/01/2026 - 10/05/2026) |
| Dot DATN | [Dien dot DATN] |

### 1. Muc tieu va yeu cau cua DATN

- Xay dung he thong quan ly trung tam day them theo mo hinh single-tenant, ho tro quan ly giao vien, hoc sinh, lop hoc, phong hoc, chuong trinh dao tao va khoa hoc.
- Tich hop module xep lich tu dong theo huong bai toan CSP, co kha nang tao preview, kiem tra rang buoc cung va ho tro commit ket qua.
- Tich hop quy trinh AI audit tai lieu giang day, gom upload tai lieu, OCR/Gemini o muc stub demo, gan nhan canh bao va phe duyet boi can bo kiem duyet.
- Xay dung giao dien web cho cac vai tro Admin, Teacher, Compliance Officer; dam bao co the demo cac luong chinh tren frontend va backend.
- He thong phai bam theo kien truc Clean Architecture, REST API, co tai lieu mo hinh hoa va huong toi tuan thu Thong tu 29/2024/TT-BGDDT.

### 2. Pham vi, noi dung cong viec va san pham can nop

- Khao sat bai toan, phan tich nghiep vu va mo hinh hoa he thong bang Use Case, Class Diagram, Sequence Diagram va tai lieu CSP.
- Thiet ke va xay dung backend bang Golang voi Gin, GORM, Wire, PostgreSQL; to chuc ma nguon theo Clean Architecture.
- Phat trien cac module chinh: xac thuc nguoi dung, CRUD giao vien, hoc sinh, lop hoc, phong hoc, chuong trinh, khoa hoc.
- Phat trien module xep lich tu dong CSP gom cau truc du lieu, backtracking, MRV, forward checking, hard constraints va API preview/commit.
- Phat trien module AI audit tai lieu gom upload, trich xuat text o muc stub, phan tich AI o muc stub, gan nhan, luu audit log va phe duyet compliance.
- Xay dung frontend ReactJS + TypeScript cho cac man hinh dang nhap, quan tri, quan ly lop/phong/hoc sinh/giao vien, xep lich va kiem duyet tai lieu.
- San pham can nop: ma nguon backend/frontend, CSDL va migration, tai lieu mo hinh hoa, mo ta API, tai lieu huong dan setup/chay va bao cao DATN.

### 3. Du lieu dau vao, gia thiet, tieu chuan/quy chuan, phan mem/cong cu su dung (neu co)

- Du lieu dau vao: thong tin giao vien, hoc sinh, lop hoc, phong hoc, chuong trinh, khoa hoc, lich hoc, tai lieu giang day, ket qua audit va thong tin nguoi dung.
- Gia thiet: he thong trien khai phuc vu demo va nghien cuu trong pham vi mot trung tam; module OCR/Gemini hien o muc stub/scaffold de hoan thien luong xu ly va danh gia mo hinh.
- Tieu chuan/quy chuan ap dung: Thong tu 29/2024/TT-BGDDT; nguyen tac RESTful API; kien truc phan lop Clean Architecture.
- Phan mem/cong cu su dung: Golang, Gin, GORM, Google Wire, PostgreSQL, ReactJS, TypeScript, Vite, MUI, Redux Toolkit, Swagger, Docker, Git.

### 4. Ke hoach/moc tien do du kien

| STT | Ke hoach | Thoi gian (moc) | Ghi chu |
|---|---|---|---|
| 1 | Khao sat de tai va xac dinh pham vi thuc hien.<br>- Tim hieu bai toan quan ly trung tam day them.<br>- Xac dinh muc tieu san pham, doi tuong su dung va cac chuc nang chinh.<br>- Lap danh sach dau vao, dau ra va san pham du kien cua DATN. | Tuan 1 (26/01/2026 - 01/02/2026) | Can chot som pham vi de tranh mo rong de tai. |
| 2 | Nghien cuu co so nghiep vu va quy dinh lien quan.<br>- Doi chieu Thong tu 29/2024/TT-BGDDT de rut ra cac rang buoc can ap dung.<br>- Xac dinh cac vai tro Admin, Teacher, Student/Parent, Compliance Officer.<br>- Liet ke cac luong nghiep vu tong quan can uu tien cho demo. | Tuan 2 (02/02/2026 - 08/02/2026) | Luu y rang buoc ve khung gio hoc va kiem soat noi dung giang day. |
| 3 | Phan tich he thong va mo hinh hoa nghiep vu tong the.<br>- Ve use case tong quan va use case theo vai tro.<br>- Xay dung cac luong nghiep vu chinh: dang nhap, quan ly lop hoc, scheduling, audit tai lieu.<br>- Xac dinh yeu cau chuc nang va phi chuc nang cua he thong. | Tuan 3 (09/02/2026 - 15/02/2026) | Nen chot use case truoc khi vao thiet ke chi tiet. |
| 4 | Thiet ke kien truc va mo hinh du lieu nen tang.<br>- Lua chon Clean Architecture cho backend va cau truc frontend.<br>- Thiet ke entity chinh, moi quan he du lieu, class diagram va sequence diagram co ban.<br>- Xay dung khung project, tai lieu setup va quy trinh phat trien. | Tuan 4 (16/02/2026 - 22/02/2026) | Can thong nhat ten bang, ten entity va quy uoc API ngay tu dau. |
| 5 | Trien khai nen tang backend va module xac thuc nguoi dung.<br>- Cau hinh project Golang, Gin, GORM, Wire, migration va logger.<br>- Xay dung cac chuc nang login, logout, refresh token, forgot/reset password, OTP verify.<br>- Tich hop xu ly loi va cau truc response dung chung. | Tuan 5 (23/02/2026 - 01/03/2026) | API auth la nen tang cho cac module sau, can on dinh som. |
| 6 | Trien khai cac module quan ly dao tao cot loi o backend.<br>- Xay dung CRUD giao vien, hoc sinh, phong hoc, lop hoc.<br>- Bo sung cac xu ly gan giao vien, them/xoa hoc sinh khoi lop, xem roster.<br>- Hoan thien migration cho cac bang nghiep vu chinh. | Tuan 6 (02/03/2026 - 08/03/2026) | Can kiem tra rang buoc khoa ngoai va du lieu mau de test API. |
| 7 | Mo rong backend cho chuong trinh dao tao va khoa hoc.<br>- Xay dung CRUD Program va Course.<br>- Thiet ke lien ket Program - Course.<br>- Hoan thien cac endpoint phuc vu trang quan tri hoc vu. | Tuan 7 (09/03/2026 - 15/03/2026) | Nen giu dong nhat cau truc request/response giua cac module CRUD. |
| 8 | Nghien cuu va tai lieu hoa bai toan xep lich theo CSP.<br>- Xac dinh variables, domains, hard constraints, soft constraints.<br>- Xay dung class diagram, sequence diagram va pseudocode solver.<br>- Chot huong trien khai backtracking, MRV, forward checking. | Tuan 8 (16/03/2026 - 22/03/2026) | Can bam sat cac rang buoc thuc te cua trung tam va TT29. |
| 9 | Trien khai module scheduling o backend.<br>- Cai dat cau truc du lieu scheduling.<br>- Trien khai backtracking, MRV, forward checking va hard constraints.<br>- Xay dung API preview/commit va co che luu ket qua de demo. | Tuan 9 (23/03/2026 - 29/03/2026) | Soft constraints co the thuc hien o muc co ban neu khong kip toi uu sau. |
| 10 | Trien khai module AI audit tai lieu o backend.<br>- Xay dung luong upload tai lieu, trich xuat text, goi dich vu AI o muc stub.<br>- Luu audit log, label va quyet dinh phe duyet.<br>- Hoan thien API cho giao vien tai tai lieu va can bo kiem duyet danh gia. | Tuan 10 (30/03/2026 - 05/04/2026) | Can ghi ro phan nao la stub, phan nao da luu duoc vao CSDL. |
| 11 | Phat trien frontend cho cac module quan tri co ban.<br>- Khoi tao giao dien dang nhap va phan quyen truy cap.<br>- Xay dung cac trang quan ly giao vien, hoc sinh, lop hoc, phong hoc, chuong trinh va khoa hoc.<br>- Ket noi frontend voi cac API backend da hoan thanh. | Tuan 11 (06/04/2026 - 12/04/2026) | Uu tien cac man hinh co the thao tac duoc du lieu that de demo. |
| 12 | Phat trien frontend cho scheduling va AI audit.<br>- Xay dung trang trigger xep lich, xem preview ket qua va hien thi xung dot.<br>- Xay dung trang tai lieu giao vien, xem ket qua AI audit.<br>- Xay dung hang doi compliance va dialog phe duyet/tu choi tai lieu. | Tuan 12 (13/04/2026 - 19/04/2026) | Can chuan bi du lieu mau de demo duoc ca 2 luong scheduling va audit. |
| 13 | Kiem thu, hoan thien tich hop va bo sung minh chung ky thuat.<br>- Kiem tra cac luong auth, CRUD, scheduling, upload/audit tai lieu.<br>- Sua loi tich hop FE-BE va bo sung du lieu test neu can.<br>- Cap nhat Swagger, anh man hinh va mo ta ket qua dat duoc. | Tuan 13 (20/04/2026 - 26/04/2026) | Uu tien sua cac loi anh huong truc tiep den demo va bao cao. |
| 14 | Hoan thien bao cao va tai lieu DATN.<br>- Tong hop qua trinh phan tich, thiet ke, trien khai va ket qua dat duoc.<br>- Viet phan danh gia, han che va huong phat trien.<br>- Chuan bi slide va kich ban demo bao ve. | Tuan 14 (27/04/2026 - 03/05/2026) | Can doi chieu giua bao cao, ma nguon va man hinh demo cho thong nhat. |
| 15 | Tong duyet va hoan tat ho so nop DATN.<br>- Chay lai toan bo cac luong demo quan trong.<br>- Chinh sua theo gop y GVHD.<br>- Hoan thien bieu mau, slide, video neu co va san sang bao ve. | Tuan 15 (04/05/2026 - 10/05/2026) | Can chuan bi phuong an demo du phong neu moi truong gap su co. |

---

## PHIEU 2 - MAU DATN-02

**NHAT KY THUC HIEN DATN**

| Muc | Noi dung de dien |
|---|---|
| Ho va ten SV | Nguyen The Ha |
| Ma SV | 61165 |
| Ten de tai | Xay dung he thong quan ly trung tam day them tich hop AI ho tro kiem soat chat luong dao tao |
| Dot | [Dien dot] |
| GVHD | [Dien ten GVHD] |
| Don vi | [Dien don vi GVHD] |
| GV dong huong dan (neu co) | [De trong neu khong co] |
| Don vi | [Dien neu co] |

### Tuan 1: 26/01/2026 - 01/02/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Khoi tao project backend, to chuc cau truc thu muc, nghien cuu kien truc va quy trinh phat trien | Da khoi tao repo Go monorepo, bo khung `cmd/internal/pkg`, bo tai lieu setup va Makefile co ban | Clean Architecture, Golang project structure, Git, Makefile | Can tiep tuc hoan thien migration va quy trinh generate | Nam duoc khung tong the de tai, can tach ro hon cac lop xu ly ngay tu dau |
| Lam viec theo nhom | 10 gio | Trao doi huong de tai va pham vi chuc nang can demo | Xac dinh de tai huong quan ly trung tam day them, tich hop scheduling va AI audit | Phan tich yeu cau, dac ta bai toan, tong hop y kien | Can chot tiep cac vai tro nguoi dung va ranh gioi scope | Rut ra can uu tien luong demo som de dinh huong kien truc on dinh |

### Tuan 2: 02/02/2026 - 08/02/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Xay dung module xac thuc va chuc nang nguoi dung | Da thuc hien login, logout, refresh token, forgot/reset password, OTP verify o backend; khoi tao frontend React | JWT, REST API, xu ly loi, React + Vite | Can dong bo flow auth giua backend va frontend, bo tri role route ro hon | Hoan thanh nen tang xac thuc, nhung can test ky cac edge case ve token |
| Lam viec theo nhom | 10 gio | Ra soat use case nguoi dung va bo tri man hinh co ban | Da xac dinh vai tro Admin, Teacher, Student/Parent, Compliance Officer; thong nhat huong giao dien co phan quyen | Phan tich use case, UI flow, role-based access | Can cu the hoa dashboard theo tung vai tro | Kinh nghiem la nen neo giao dien vao use case truoc khi code chi tiet |

### Tuan 3: 09/02/2026 - 15/02/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Xay dung module quan ly giao vien va giao dien ban dau | Da hoan thanh CRUD giao vien, tim kiem/loc, lay thong tin ho so va mot so thanh phan UI lien quan | Gin, GORM, repository/usecase, React component | Can bo sung test va chuan hoa response FE/BE | Tien do tot, can tiep tuc giu mot contract API on dinh cho frontend |
| Lam viec theo nhom | 10 gio | Mo hinh hoa nghiep vu chi tiet theo vai tro | Da bo sung use case cho Admin, Teacher, Student/Parent, Compliance Officer va AI Agent trong tai lieu | UML, use case modeling, phan tich nghiep vu | Can lien ket chat hon giua use case va danh sach API/DB | Rut kinh nghiem nen song song hoa tai lieu va code de giam sai lech |

### Tuan 4: 16/02/2026 - 22/02/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Tai lieu hoa bai toan xep lich va thiet ke he thong | Da hoan thanh tai lieu CSP, class diagram core/AI audit/scheduling, sequence diagram login, teacher CRUD, scheduling va audit | CSP modeling, UML class/sequence, phan tich rang buoc | Can bo sung ERD tong the va mo ta rang buoc du lieu chi tiet hon | Phan mo hinh hoa giup dinh hinh ro hon cach trien khai module scheduling |
| Lam viec theo nhom | 10 gio | Ra soat pham vi tuan thu va gia tri khoa hoc cua de tai | Da xac dinh ro 2 diem nhan: scheduling theo CSP va AI ho tro kiem duyet tai lieu, gan voi Thong tu 29/2024/TT-BGDDT | Doi chieu quy dinh, tong hop yeu cau compliance | Can cu the hoa them bo test cho cac rang buoc TT29 | Kinh nghiem la moi tinh nang can gan voi mot gia tri hoc thuat va thuc tien ro rang |

### Tuan 5: 23/02/2026 - 01/03/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Hoan thien cac module quan ly dao tao cot loi | Da hoan thanh CRUD room, class, student, program, course; bo sung usecase them/xoa hoc sinh khoi lop va gan giao vien | CRUD design, migration, API design, normalize du lieu | Can bo sung cac tinh nang ngoai demo nhu diem danh, hoc phi, ket qua hoc tap | Khoi core module da ro rang hon, can uu tien tinh nang co the demo tron ven |
| Lam viec theo nhom | 10 gio | Thong nhat danh sach tinh nang uu tien cho demo | Da chot nhom tinh nang uu tien: class roster, scheduling preview, AI audit, giao dien admin/compliance/teacher | Quan ly task, uu tien backlog, doi chieu pham vi | Can cat bot cac muc chua can thiet cho ky demo | Rut kinh nghiem la can uu tien depth over breadth de tranh do dang nhieu module |

### Tuan 6: 02/03/2026 - 08/03/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Trien khai scheduling CSP va AI audit o muc demo, dong thoi hoan thien frontend lien quan | Da co cau truc du lieu CSP, backtracking, MRV, forward checking, hard constraints, API preview/commit; da co upload tai lieu, AI audit stub, gan nhan, queue compliance, man hinh teacher/compliance va scheduling | Thuat toan CSP, heuristics, React state, tich hop API, luu audit log | Chua hoan thien soft constraints, OCR/Gemini that, preview/download tai lieu va test hieu nang | Dat duoc bo khung demo ro rang, nhung can thang than ghi chu cac phan moi o muc scaffold/stub |
| Lam viec theo nhom | 10 gio | Tong hop tien do, kiem tra luong demo va danh gia cac phan con thieu | Da tong hop Top 10 uu tien demo, bo sung class roster UI, scheduling preview FE, AI audit BE/FE va cap nhat tai lieu tien do | Demo review, doi chieu task, giao tiep ky thuat | Can tiep tuc bo sung kiem thu tich hop, tai lieu bao cao va kich ban bao ve | Rut ra can ghi ro muc do hoan thanh cua tung module de tranh mo ta vuot qua hien trang |

### Tuan 7: 09/03/2026 - 15/03/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Hoan thien cac luong demo tren frontend va ket noi voi backend da co | Da bo sung giao dien class roster, scheduling preview, teacher documents, compliance queue va cac dialog chi tiet de phuc vu demo thong suot | React, TypeScript, tich hop API, xu ly state va luong nguoi dung | Can tiep tuc dong bo response FE-BE, bo sung preview/download tai lieu va bo test giao dien | Kinh nghiem la can uu tien luong demo co the trinh bay tron ven truoc khi mo rong them chuc nang |
| Lam viec theo nhom | 10 gio | Rasoat backlog demo, chot cac hang muc uu tien va cach kiem tra nhanh | Da tong hop Top 10 uu tien demo, doi chieu tung module voi endpoint va man hinh su dung thuc te | Quan ly task, review demo, phan chia muc uu tien | Can tiep tuc bo sung kich ban test tich hop va tai lieu minh chung cho moi luong | Rut ra can mo ta ro phan nao da hoan thanh, phan nao moi o muc scaffold de tranh danh gia sai |

### Tuan 8: 16/03/2026 - 22/03/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Ra soat cac module da co, cap nhat tai lieu va chuan bi noi dung bao cao | Da doi chieu ma nguon voi danh sach task, xac dinh cac module da hoan thanh va cac phan chua bat dau nhu testing, deployment, dashboard thong ke, chatbot | Tong hop tai lieu, danh gia tien do, viet bao cao ky thuat | Can bo sung bang doi chieu giua yeu cau - chuc nang - minh chung | Can tiep tuc giam khoang cach giua implementation va phan bao cao |
| Lam viec theo nhom | 10 gio | Trao doi cach trinh bay ket qua DATN theo mau quy dinh | Da xac dinh cach dien DATN-01 theo muc tieu/pham vi/tien do va DATN-02 theo nhat ky tuan | Ky nang tong hop, chuan hoa tai lieu, trinh bay hoc thuat | Can xin xac nhan them thong tin hanh chinh tu GVHD | Kinh nghiem la nen chot som mau bieu de cap nhat dinh ky tu dau ky |

### Tuan 9: 23/03/2026 - 29/03/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Tong hop hien trang project, kiem tra muc do hoan thanh cua cac luong demo va cac phan con thieu | Da tong hop duoc hien trang: da xong khung he thong va cac luong demo chinh; testing, toi uu, deploy, dashboard thong ke va mot so tinh nang hoc vu van chua hoan tat | Tong hop minh chung, doi chieu tien do voi ke hoach | Can tiep tuc kiem thu chuc nang, viet bao cao va chuan bi demo/bao ve | Danh gia chung: project dat muc prototype/demo tot, can bo sung phan kiem thu va hoan thien de tang tinh day du |
| Lam viec theo nhom | 10 gio | Ra soat noi dung ho so DATN va thong nhat cach the hien tien do theo tuan | Da chuan bi duoc bo noi dung co so de dien cac bieu mau tien do va nhat ky theo dung tien do thuc te cua project | Chuan hoa noi dung, tong hop minh chung, giao tiep ky thuat | Can thay the them thong tin ca nhan, GVHD, don vi va chu ky truoc khi nop | Rut kinh nghiem can cap nhat nhat ky hang tuan ngay khi hoan thanh cong viec de tranh phai tong hop lai vao cuoi dot |

---

## Ghi chu de dien nhanh vao mau

- Neu can viet gon hon o phieu DATN-01, uu tien giu nguyen ten de tai, 4 muc lon va bang tien do 6-7 moc.
- O phieu DATN-02, co the chep moi tuan 2 dong: `Lam viec doc lap` va `Lam viec theo nhom`.
- Cac muc nen giu cach dien than trong:
  - Scheduling: ghi `muc scaffold/demo` vi soft constraints va toi uu chua hoan tat.
  - AI audit: ghi `stub OCR/Gemini` neu GVHD yeu cau mo ta dung hien trang.
  - Testing/deployment: ghi `dang thuc hien` hoac `chua hoan tat`, khong nen danh dau da xong.
