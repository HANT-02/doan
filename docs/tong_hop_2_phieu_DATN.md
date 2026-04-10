# Tong hop 2 phieu DATN theo noi dung do an hien tai

Tai lieu nay duoc tong hop tu:
- `/Users/hant/golang/doan/PROJECT_TASKS.md`
- `/Users/hant/golang/doan/docs/ke_hoach_phan_hoi_gvhd_2026-04-09.md`
- ma nguon backend/frontend hien tai trong repo
- tien do trien khai thuc te den ngay `10/04/2026`

Luu y:
- Cac truong hanh chinh khong xac dinh chac chan tu repo duoc de dang `[Dien ...]`.
- Thong tin sinh vien suy ra tu tai lieu cu trong repo: `Nguyen The Ha - 61165 - CS2`.
- Noi dung duoi day da duoc sua lai theo huong do an moi:
  - bo nhanh `AI Audit/Compliance`,
  - lay scheduling lam trong tam ky thuat,
  - bo sung `Shift`,
  - benchmark 3 solver `Graph Coloring`, `CP-SAT`, `Tabu Search`,
  - bo sung huong `Predictive Analytics` cho bai toan `AT_RISK classification`.

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
| Ten de tai DATN | Xay dung he thong quan ly trung tam day them, trong tam la xep lich thong minh va du bao sinh vien co nguy co hoc kem |
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

- Xay dung he thong quan ly trung tam day them theo mo hinh single-tenant, ho tro cac nghiep vu cot loi: quan ly giao vien, hoc sinh, lop hoc, phong hoc, chuong trinh dao tao, khoa hoc va tai khoan.
- Nang module xep lich thanh diem nhan ky thuat cua do an, trong do du lieu thoi gian duoc chuan hoa bang thuc the `Shift` thay cho cach su dung khung gio roi rac.
- Nghien cuu, cai dat va so sanh 3 thuat toan xep lich gom `Graph Coloring + heuristic`, `CP-SAT` va `Tabu Search`, tu do benchmark va lua chon solver phu hop nhat cho he thong chinh.
- Xay dung co che scheduling co preview, kiem tra hard constraints, soft score, benchmark admin API noi bo va commit ket qua xuong `lesson`.
- Bo sung bai toan `Predictive Analytics` theo huong machine learning de du bao sinh vien co nguy co hoc kem (`AT_RISK classification`) dua tren du lieu hoc tap va van hanh.
- Dam bao he thong bam theo Clean Architecture, REST API va phu hop dinh huong nghien cuu/ung dung cua do an tot nghiep.

### 2. Pham vi, noi dung cong viec va san pham can nop

- Khao sat bai toan quan ly trung tam day them, phan tich nghiep vu va mo hinh hoa he thong bang Use Case, Class Diagram, Sequence Diagram va tai lieu scheduling.
- Thiet ke va xay dung backend bang Golang voi Gin, GORM, Wire, PostgreSQL; to chuc ma nguon theo Clean Architecture.
- Phat trien cac module quan ly cot loi: giao vien, hoc sinh, lop hoc, phong hoc, chuong trinh dao tao, khoa hoc, tai khoan.
- Xay dung module `Shift` de quan ly ca hoc va refactor `class_schedule` sang `shift_id`.
- Phat trien module scheduling gom preview, commit, conflict messaging, hard constraints, soft score va benchmark 3 solver qua admin API noi bo.
- Cai dat 3 solver scheduling theo abstraction chung o tang service:
  - `Graph Coloring + heuristic`,
  - `CP-SAT`,
  - `Tabu Search`.
- Xay dung huong `Predictive Analytics` cho bai toan phan lop sinh vien `AT_RISK / NOT_AT_RISK`, gom xac dinh feature, huan luyen, danh gia va endpoint du bao.
- Xay dung frontend ReactJS + TypeScript cho cac man hinh quan tri va scheduling, bao dam demo duoc luong nghiep vu chinh.
- San pham can nop: ma nguon backend/frontend, migration CSDL, tai lieu mo hinh hoa, tai lieu benchmark scheduling, mo ta huong predictive analytics, slide va bao cao DATN.

### 3. Du lieu dau vao, gia thiet, tieu chuan/quy chuan, phan mem/cong cu su dung (neu co)

- Du lieu dau vao:
  - thong tin giao vien, hoc sinh, lop hoc, phong hoc,
  - chuong trinh dao tao, khoa hoc,
  - ca hoc `Shift`,
  - lich mau cua lop `class_schedule`,
  - du lieu `lesson` de commit thoi khoa bieu,
  - du lieu hoc tap/phu tro de xay dung bai toan du bao `AT_RISK`.
- Gia thiet:
  - he thong trien khai phuc vu nghien cuu va demo trong pham vi mot trung tam day them,
  - benchmark scheduling duoc chay tren bo du lieu noi bo cua he thong,
  - predictive analytics duoc huan luyen bang backend hien tai va du lieu duoc tong hop trong pham vi do an.
- Tieu chuan/quy chuan ap dung:
  - Thong tu 29/2024/TT-BGDDT,
  - nguyen tac RESTful API,
  - Clean Architecture,
  - cac nguyen tac danh gia mo hinh va benchmark scheduling.
- Phan mem/cong cu su dung:
  - Golang, Gin, GORM, Google Wire,
  - PostgreSQL,
  - ReactJS, TypeScript, Vite, MUI, Redux Toolkit,
  - Swagger,
  - Git, Docker,
  - cac thu vien/ky thuat phuc vu benchmark va machine learning trong backend hien tai.

### 4. Ke hoach/moc tien do du kien

| STT | Ke hoach | Thoi gian (moc) | Ghi chu |
|---|---|---|---|
| 1 | Khao sat de tai va xac dinh pham vi nghien cuu.<br>- Tim hieu bai toan quan ly trung tam day them.<br>- Xac dinh nhom chuc nang cot loi cua he thong.<br>- Chot huong scheduling la diem nhan ky thuat cua DATN. | Tuan 1 (26/01/2026 - 01/02/2026) | Can chot som muc tieu de tai de tranh mo rong qua nhieu nhanh khong can thiet. |
| 2 | Nghien cuu nghiep vu va rang buoc bai toan.<br>- Doi chieu quy dinh va cac rang buoc van hanh hoc tap.<br>- Xac dinh cac actor va use case tong quan.<br>- Liet ke dau vao, dau ra va cac rang buoc scheduling. | Tuan 2 (02/02/2026 - 08/02/2026) | Nen xac dinh ro hard constraints va soft constraints ngay tu dau. |
| 3 | Thiet ke kien truc va khoi tao he thong.<br>- Lua chon Clean Architecture cho backend va cau truc frontend.<br>- Khoi tao project Golang/React, migration, logger, auth flow.<br>- Dat nen tang cho cac module quan ly cot loi. | Tuan 3 (09/02/2026 - 15/02/2026) | Can thong nhat quy uoc API, ten entity va response dung chung. |
| 4 | Trien khai cac module CRUD cot loi o backend va frontend.<br>- Quan ly giao vien, hoc sinh, lop hoc, phong hoc.<br>- Quan ly chuong trinh dao tao va khoa hoc.<br>- Hoan thien cac luong roster va gan giao vien. | Tuan 4 (16/02/2026 - 22/02/2026) | Uu tien cac module su dung truc tiep cho scheduling ve sau. |
| 5 | Mo hinh hoa va xay dung scheduling baseline.<br>- Xac dinh variables, domains, hard constraints, soft score.<br>- Xay dung preview/commit scheduling ban dau.<br>- Tai lieu hoa logic scheduling va conflict messaging. | Tuan 5 (23/02/2026 - 01/03/2026) | Giai doan nay tao bo khung de ve sau nang cap bang nhieu solver. |
| 6 | Hoan thien giao dien va tich hop scheduling baseline.<br>- Xay dung trang scheduling preview tren frontend.<br>- Kiem tra luong FE-BE va sua contract preview/commit.<br>- Tong hop cac van de UX va du lieu nghiep vu. | Tuan 6 (02/03/2026 - 08/03/2026) | Can on dinh luong preview truoc khi mo rong benchmark. |
| 7 | Tiep nhan gop y GVHD va dieu chinh scope do an.<br>- Loai bo nhanh AI Audit/Compliance khoi backlog chinh.<br>- Chot huong moi: `Shift + 3 solver scheduling + benchmark + AT_RISK classification`.<br>- Cap nhat muc tieu, backlog va tai lieu huong nghien cuu. | Tuan 7 (09/03/2026 - 15/03/2026) | Can giu scope du sau nhung van kha thi trong thoi gian DATN. |
| 8 | Xay dung module `Shift` va chuan hoa du lieu thoi gian cho scheduling.<br>- Thiet ke schema `shifts` va CRUD admin.<br>- Them UI quan ly ca hoc.<br>- Chuan bi refactor `class_schedule` sang `shift_id`. | Tuan 8 (16/03/2026 - 22/03/2026) | `Shift` la du lieu nen, can on dinh truoc khi benchmark solver. |
| 9 | Refactor scheduling sang `Shift`.<br>- Chuyen `class_schedule` sang `shift_id`.<br>- Sua domain generation, preview output va conflict message theo `Shift`.<br>- Kiem tra backward impact voi du lieu hien co. | Tuan 9 (23/03/2026 - 29/03/2026) | Can co migration backfill de khong mat du lieu lich mau cu. |
| 10 | Cai dat 3 solver scheduling theo abstraction chung.<br>- Implement `Graph Coloring + heuristic`.<br>- Implement `CP-SAT`.<br>- Implement `Tabu Search`.<br>- Dung chung scorer, constraint checker va output model. | Tuan 10 (30/03/2026 - 05/04/2026) | Solver duoc to chuc o tang service, use case chi dung abstraction. |
| 11 | Xay dung benchmark scheduling va danh gia solver.<br>- Tao admin API benchmark cho 3 solver.<br>- Chay cung mot input de lay metric feasibility, hard violations, soft score, runtime.<br>- Chuan bi bang so sanh va co so lua chon solver chinh. | Tuan 11 (06/04/2026 - 12/04/2026) | Day la phan tao gia tri nghien cuu ro nhat cho do an. |
| 12 | Chon solver chinh va gan vao scheduling API production-like.<br>- Lua chon solver tot nhat sau benchmark.<br>- Inject vao use case scheduling chinh.<br>- Kiem tra lai luong preview, commit va UI scheduling. | Tuan 12 (13/04/2026 - 19/04/2026) | Benchmark van giu thanh admin API noi bo, khong dua cho nguoi dung cuoi. |
| 13 | Xac dinh bai toan predictive analytics `AT_RISK classification`.<br>- Chot nhan `AT_RISK / NOT_AT_RISK`.<br>- Xac dinh nguon du lieu, feature set va metric danh gia.<br>- Chuan bi pipeline train/evaluate trong backend hien tai. | Tuan 13 (20/04/2026 - 26/04/2026) | Chi chon mot bai toan du bao de giu do sau cho do an. |
| 14 | Trien khai pipeline machine learning va endpoint du bao.<br>- Huan luyen, danh gia mo hinh classification.<br>- Tao API du bao va metadata mo hinh.<br>- Chuan bi giao dien canh bao sinh vien nguy co hoc kem neu kip. | Tuan 14 (27/04/2026 - 03/05/2026) | Can uu tien metric va minh chung hoc thuat hon la mo rong qua nhieu UI. |
| 15 | Kiem thu, tong hop bao cao va tong duyet.<br>- Kiem tra end-to-end core modules, scheduling va predictive analytics.<br>- Tong hop bang benchmark, ket qua danh gia mo hinh va hinh anh demo.<br>- Chinh sua theo gop y GVHD, hoan tat bieu mau, slide va bao cao. | Tuan 15 (04/05/2026 - 10/05/2026) | Can doi chieu giua ma nguon, so lieu benchmark va noi dung bao cao de thong nhat. |

---

## PHIEU 2 - MAU DATN-02

**NHAT KY THUC HIEN DATN**

| Muc | Noi dung de dien |
|---|---|
| Ho va ten SV | Nguyen The Ha |
| Ma SV | 61165 |
| Ten de tai | Xay dung he thong quan ly trung tam day them, trong tam la xep lich thong minh va du bao sinh vien co nguy co hoc kem |
| Dot | [Dien dot] |
| GVHD | [Dien ten GVHD] |
| Don vi | [Dien don vi GVHD] |
| GV dong huong dan (neu co) | [De trong neu khong co] |
| Don vi | [Dien neu co] |

### Tuan 1: 26/01/2026 - 01/02/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Khao sat de tai, xac dinh bai toan va huong nghien cuu tong quan | Da xac dinh de tai huong quan ly trung tam day them, trong do xep lich la bai toan ky thuat trung tam can dau tu sau nay | Phan tich yeu cau, tong hop bai toan, chot scope ban dau | Can doi chieu them voi yeu cau hoc thuat de giu duoc gia tri nghien cuu | Da co dinh huong tong quan, can tiep tuc chot diem nhan ky thuat ro hon |
| Lam viec theo nhom | 10 gio | Trao doi pham vi module va san pham du kien | Da thong nhat nhom module cot loi can co de demo: giao vien, hoc sinh, lop hoc, phong hoc, chuong trinh, khoa hoc | Giao tiep ky thuat, phan chia backlog | Can uu tien module lien quan truc tiep den scheduling | Rut kinh nghiem can phan loai som phan nao la cot loi, phan nao la mo rong |

### Tuan 2: 02/02/2026 - 08/02/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Nghien cuu nghiep vu, actor va cac rang buoc van hanh | Da xac dinh actor chinh va cac rang buoc co anh huong den scheduling nhu lop, giao vien, phong hoc va khung hoc | Use case modeling, phan tich nghiep vu | Can chuyen tiep cac rang buoc nay thanh mo hinh du lieu va API | Da hinh dung duoc bai toan scheduling khong chi la CRUD ma la bai toan rang buoc |
| Lam viec theo nhom | 10 gio | Doi chieu scope ban dau voi gia tri bao cao DATN | Da thong nhat huong can co module scheduling duoc mo ta ro ve thuat toan, khong chi dung o muc giao dien | Phan tich hoc thuat, tong hop yeu cau GVHD | Can bo sung tai lieu mo hinh hoa ngay tu giai doan dau | Kinh nghiem la can gan moi tinh nang voi gia tri hoc thuat cu the |

### Tuan 3: 09/02/2026 - 15/02/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Khoi tao kien truc he thong, auth flow va cau truc project | Da khoi tao backend Golang theo Clean Architecture, migration, logger, auth flow va frontend React ban dau | Golang, Gin, GORM, Wire, React, TypeScript | Can tiep tuc on dinh auth contract giua FE-BE | Da dat nen tang tot, can giu tinh dong nhat kien truc khi mo rong module |
| Lam viec theo nhom | 10 gio | Ra soat ten entity, quy uoc API va cau truc response | Da thong nhat quy uoc dat ten va cach to chuc use case/repository cho cac module CRUD | API design, clean code, organization | Can bo sung them tai lieu quy uoc de giam sai lech giua frontend va backend | Kinh nghiem la thong nhat quy uoc som giup giam loi tich hop ve sau |

### Tuan 4: 16/02/2026 - 22/02/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Trien khai cac module CRUD cot loi cua he thong | Da hoan thien phan lon CRUD giao vien, hoc sinh, lop hoc, phong hoc, chuong trinh va khoa hoc | CRUD design, migration, use case/repository pattern | Can tiep tuc kiem tra roster, gan giao vien va contract detail API | Core module da ro rang hon, tao duoc nen du lieu cho scheduling |
| Lam viec theo nhom | 10 gio | Kiem tra kha nang demo cac module quan tri | Da danh gia duoc nhom module co the demo duoc tren UI va xac dinh nhung cho can sua UX de phuc vu trinh bay | Demo review, quan ly task, UI flow | Can tiep tuc sua cac loi giao dien va flow thao tac admin | Rut kinh nghiem can lam song song core API va man hinh thao tac that |

### Tuan 5: 23/02/2026 - 01/03/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Mo hinh hoa va cai dat scheduling baseline | Da xay dung preview/commit scheduling ban dau, co hard constraints, conflict messaging va tai lieu hoa logic scheduling hien tai | Scheduling modeling, backtracking, heuristic, conflict handling | Chua co benchmark, chua co nhieu solver va du lieu thoi gian chua chuan hoa | Baseline da chay duoc, nhung can nang cap de tao gia tri nghien cuu ro hon |
| Lam viec theo nhom | 10 gio | Tong hop bai toan scheduling can mo rong tiep | Da nhan dien cac han che: du lieu thoi gian con roi rac, solver moi o muc mot huong, thieu benchmark va tieu chi so sanh | Phan tich thuat toan, tong hop han che | Can bo sung thuc the thoi gian chuan va nhieu solver de dat yeu cau GVHD | Kinh nghiem la can nhin scheduling nhu mot bai toan nghien cuu doc lap |

### Tuan 6: 02/03/2026 - 08/03/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Hoan thien frontend admin va scheduling baseline | Da ket noi frontend voi cac API chinh, sua cac luong quan ly giao vien, lop hoc, phong hoc, chuong trinh/khoa hoc va man scheduling preview | React, TypeScript, RTK Query, sua UX/contract FE-BE | Can tiep tuc chuan hoa preview scheduling va fix nhung mismatch con lai | Giao dien da duoc su dung de kiem chung lai thiet ke backend rat hieu qua |
| Lam viec theo nhom | 10 gio | Ra soat toan bo luong demo chinh cua he thong | Da xac dinh duoc nhom luong se trinh bay trong bao cao: core admin modules va scheduling | Tong hop minh chung, scenario test | Can tiep tuc lam day scheduling truoc khi mo sang nhanh khac | Rut kinh nghiem la can tap trung mot diem nhan ky thuat thay vi dan trai qua nhieu nhanh |

### Tuan 7: 09/03/2026 - 15/03/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Tiep nhan gop y GVHD va dieu chinh scope do an | Da bo nhanh `AI Audit/Compliance` khoi backlog chinh va chot scope moi: `Shift + 3 solver scheduling + benchmark + AT_RISK classification` | Tong hop phan hoi GVHD, danh gia lai pham vi | Can cap nhat tat ca task, tai lieu va bieu mau theo scope moi | Viec cat bo dung luc giup de tai sau hon va thong nhat hon |
| Lam viec theo nhom | 10 gio | Lap lai ke hoach trien khai sau dieu chinh | Da viet ke hoach chi tiet cho scheduling benchmark, module `Shift` va predictive analytics | Planning, backlog refinement, architecture review | Can chuyen nhanh ke hoach thanh task ky thuat cu the | Rut kinh nghiem la phai khoa quyet dinh kien truc truoc khi code tiep |

### Tuan 8: 16/03/2026 - 22/03/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Trien khai module `Shift` cho quan ly ca hoc | Da thiet ke schema, migration, entity, repository, use case, controller va UI CRUD cho `Shift` | Thiet ke du lieu, CRUD, migration, UI admin | Can noi `Shift` vao du lieu lich mau va scheduling domain | `Shift` da tro thanh don vi thoi gian chuan, tao co so tot cho benchmark sau nay |
| Lam viec theo nhom | 10 gio | Kiem tra kha nang su dung `Shift` trong he thong | Da danh gia cac diem can refactor, dac biet la `class_schedule` va scheduling preview | Impact analysis, refactor planning | Can migration backfill du lieu cu de dam bao chay duoc tren DB hien tai | Kinh nghiem la nen lam du lieu nen truoc roi moi doi solver |

### Tuan 9: 23/03/2026 - 29/03/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Refactor scheduling sang `Shift` | Da chuyen `class_schedule` sang `shift_id`, sua domain generation, preview output va conflict messaging theo `Shift` | Refactor schema, migration backfill, scheduling domain | Can tiep tuc benchmark tren 3 solver de co co so chon solver chinh | Day la moc quan trong vi scheduling khong con phu thuoc vao khung gio hardcode roi rac |
| Lam viec theo nhom | 10 gio | Kiem tra backward impact va tinh thong nhat du lieu | Da ra soat migration, preload `ClassSchedules.Shift`, preview UI va flow commit sau refactor | Review impact, FE-BE integration | Can chay benchmark va them minh chung cho bao cao | Rut kinh nghiem can sua dong bo schema, service va UI trong cung mot dot refactor |

### Tuan 10: 30/03/2026 - 05/04/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Chuan hoa kien truc scheduling de ho tro nhieu solver | Da tach `SchedulingSolver` interface, chuan hoa input/output/scorer, benchmark result model va giu use case chi phu thuoc abstraction | Clean Architecture, interface design, service abstraction | Can tiep tuc implement day du 3 solver va benchmark API thuc te | Kien truc moi giup giam lap code va tao nen benchmark cong bang hon |
| Lam viec theo nhom | 10 gio | Ra soat pattern code cho giai doan solver | Da chot rule: moi solver la mot implementation rieng o tang service, benchmark la admin API noi bo, frontend chinh khong cho chon solver | Architecture review, coding convention | Can giu sat pattern nay trong suot giai doan benchmark | Rut kinh nghiem la phai khoa kien truc truoc khi so sanh thuat toan |

### Tuan 11: 06/04/2026 - 12/04/2026

| Loai cong viec | Thoi gian lam viec | Nhiem vu duoc giao | Ket qua thuc hien | Kien thuc/ky nang ap dung | Cac van de can bo sung, chinh sua | Tu danh gia va rut kinh nghiem |
|---|---:|---|---|---|---|---|
| Lam viec doc lap | 10 gio | Cai dat 3 solver scheduling va test tren input chuan | Da implement `GraphColoringSolver`, `CPSATSolver`, `TabuSearchSolver`, bo helper dung chung va test co ban cho ca 3 solver | Graph Coloring heuristic, CP-style search, Tabu Search, reuse scorer | Chua chay benchmark admin API thuc te de lay metric so sanh va chua chon solver chinh | Da hoan thanh buoc quan trong de dua scheduling tu prototype sang huong nghien cuu co so sanh |
| Lam viec theo nhom | 10 gio | Cap nhat backlog, tai lieu va bieu mau theo hien trang moi | Da cap nhat `PROJECT_TASKS.md`, tai lieu ke hoach va tong hop lai bieu mau theo huong do an hien tai | Tong hop tai lieu, doi chieu ma nguon voi tien do | Can tiep tuc them so lieu benchmark va huong predictive analytics vao bao cao | Rut kinh nghiem la can cap nhat tai lieu song song voi implementation de bao cao khop hien trang |

---

## Ghi chu de dien nhanh vao mau

- Ten de tai nen thong nhat xuyen suot theo huong moi:
  - `Xay dung he thong quan ly trung tam day them, trong tam la xep lich thong minh va du bao sinh vien co nguy co hoc kem`.
- Trong phieu DATN-01, phan scheduling nen mo ta dung tinh chat hien tai:
  - da co `Shift`,
  - da refactor scheduling domain,
  - da implement 3 solver,
  - benchmark va chon solver chinh la buoc tiep theo.
- Trong phieu DATN-02, nen nhan manh tu tuan 7 tro di la giai doan dieu chinh scope theo gop y GVHD.
- Phan predictive analytics hien moi o muc huong nghien cuu/ke hoach trien khai, khong nen viet nhu da hoan thanh.
- Neu can viet gon hon khi chep vao mau:
  - giu lai 3 trong tam: `core modules`, `scheduling`, `predictive analytics`,
  - khong nhac lai nhanh `AI Audit/Compliance` de tranh le scope.
