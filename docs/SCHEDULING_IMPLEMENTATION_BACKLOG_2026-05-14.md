# Scheduling Implementation Backlog 2026-05-14

Tài liệu này là backlog triển khai scheduling production-like sau các vòng bổ sung nghiệp vụ ngày 2026-05-14.

Mục tiêu hiện tại:

- xếp lịch là bài toán chính để chấm điểm;
- học máy chỉ còn là hướng phát triển tương lai;
- ưu tiên thuật toán xếp lịch, bộ dữ liệu lớn, benchmark hiệu năng, độ đúng luật cứng và chất lượng luật mềm;
- chọn `CP-SAT` dựa trên số liệu thực nghiệm, không chọn bằng cảm tính;
- giao diện vận hành phải gọn, thuần tiếng Việt, ít chữ thừa, giúp giáo vụ sửa lịch nhanh.

## 1. Trạng thái hiện tại sau các đợt đã làm

Đã có trong repo:

- preview scheduling và commit preview sang buổi học;
- benchmark solver ban đầu, trong đó `CP-SAT` là solver chính;
- rule lọc lớp chưa đủ sĩ số tối thiểu trước khi đưa vào preview;
- fix rule sĩ số để đếm đúng học viên `APPROVED` và `ENROLLED`;
- rule chặn commit nếu số ca chỉnh tay vượt ngưỡng vận hành;
- lifecycle cơ bản cho buổi học: `DRAFT`, `PUBLISHED`, `HISTORY`;
- mô hình `campus`, `campus_travel_times`, `room.campus_id`;
- field kỹ năng dạng mảng trên giáo viên và khóa học;
- hard check kỹ năng trong preview;
- preview mode backend có các chế độ xếp mới, xếp lại nháp, giữ lịch đã công bố;
- API/use case dạy thay cơ bản;
- API tìm slot học bù dạng tra cứu;
- nghiệp vụ bảo lưu, hoàn tác bảo lưu và chuyển lớp;
- impact summary cơ bản: số ca đổi lịch, đổi giáo viên, đổi phòng, mức lấp đầy;
- UI teacher substitute, học bù lookup, roster bảo lưu/chuyển lớp;
- màn lịch admin, teacher, student bắt đầu dùng component lịch tuần chung.

Những điểm đã có nhưng cần chỉnh lại:

- UI preview còn rối, nhiều subtext, alert và tiếng Anh lẫn tiếng Việt;
- flow đổi slot thủ công đang liệt kê quá nhiều phương án cùng lúc;
- kéo thả đã có nhưng cần test và tinh chỉnh UX;
- lịch tuần chưa có cột khung giờ để nhìn ngang theo ca;
- preview vẫn cho chọn chế độ xếp lịch ở UI dù product direction là mặc định giữ lịch đã công bố;
- quản lý buổi học còn filter lifecycle không cần thiết với người dùng cuối;
- kỹ năng/chứng chỉ hiện là field mảng, chưa phải bảng chuẩn hóa.

## 2. Nguyên tắc thiết kế mới

Luật sản phẩm:

- không xếp lịch về quá khứ;
- nếu người dùng chọn ngày bắt đầu dự kiến trong quá khứ, hệ thống tự lấy từ ngày hiện tại hoặc ngày hợp lệ gần nhất;
- vùng quá khứ chỉ là lịch sử, không đưa vào preview và không tính conflict cho kế hoạch mới;
- mặc định preview dùng chế độ giữ lịch đã công bố;
- không hiển thị bộ lọc mode xếp mới/giữ lịch cũ trên UI vận hành;
- chỉ dùng xếp mới hoàn toàn trong benchmark hoặc debug nội bộ.

Luật UI:

- toàn bộ text hiển thị phải là tiếng Việt;
- bỏ subtext, subtitle, alert nếu không giúp thao tác nhanh hơn;
- ưu tiên chip, bảng, tooltip ngắn và trạng thái rõ;
- calendar không ghi text trống trong ô trống;
- card lịch chỉ hiện thông tin cần để ra quyết định;
- thông tin giải thích dài đưa vào hover hoặc panel phụ.

Luật chỉnh tay:

- người dùng chọn ca trước;
- sau đó chọn phòng;
- hệ thống mới hiện các vị trí hợp lệ tương ứng;
- mỗi slot/chip phải biết buổi thứ mấy của khóa học;
- tooltip slot/chip hiển thị: buổi thứ mấy, luật cứng đã qua, điểm cộng, điểm trừ, lý do điểm mềm thay đổi.

## 3. Các vấn đề cần xử lý theo nhóm

### Nhóm P0 - Lỗi đúng sai scheduling

P0.1. Không xếp lịch về quá khứ

Mục tiêu:

- preview không sinh buổi học trước ngày hiện tại;
- không tính các ngày cũ vào số slot cần xếp;
- không báo thiếu slot do lẫn dữ liệu quá khứ.

Task:

- thêm hàm chuẩn hóa `effectiveDateFrom = max(inputDateFrom, today, classStartDate nếu hợp lệ)`;
- nếu lớp đã có `start_date` trong quá khứ nhưng còn buổi chưa xếp, chỉ xếp phần từ hôm nay trở đi;
- cập nhật preview response để trả lại `effective_date_from` cho UI hiểu hệ thống đã tự điều chỉnh;
- thêm test case chọn ngày bắt đầu trong quá khứ.

Acceptance:

- chọn ngày bắt đầu trước ngày hiện tại không sinh assignment quá khứ;
- số buổi cần xếp không bị phình vì khoảng ngày cũ;
- message UI không yêu cầu người dùng tự sửa ngày nếu hệ thống có thể tự chuẩn hóa.

P0.2. Số buổi và khoảng thời gian phải theo từng lớp

Mục tiêu:

- chọn nhiều lớp hoặc chọn giáo viên dạy nhiều lớp thì mỗi lớp có số buổi, số buổi/tuần và khoảng thời gian riêng;
- không dùng chung một `date_to` suy ra từ lớp dài nhất cho mọi lớp một cách làm sai logic.

Logic mới:

- lấy `course.total_sessions` hoặc field tương đương làm số buổi cần xếp;
- lấy số buổi/tuần từ cấu hình khóa học hoặc lịch tuần lớp;
- lấy ngày bắt đầu dự kiến của lớp;
- suy ra ngày kết thúc dự kiến của từng lớp từ số buổi và số buổi/tuần;
- preview input có thể nhận một ngày bắt đầu chung, nhưng solver phải build window riêng theo lớp.

Task:

- rà model khóa học để xác định field số buổi và số buổi/tuần hiện có hoặc cần bổ sung;
- bổ sung `ClassSchedulingWindow` trong builder: `class_id`, `effective_start`, `expected_end`, `sessions_required`;
- sửa `buildVariables` và `buildDomains` để dùng window riêng của từng class;
- thêm test chọn 2 lớp có số buổi và lịch tuần khác nhau.

Acceptance:

- class A 12 buổi, class B 24 buổi không bị dùng chung session total;
- giáo viên có nhiều lớp thì preview sinh đúng số buổi từng lớp;
- UI summary từng lớp hiện đúng `đã xếp/tổng buổi`.

P0.3. Conflict với buổi đã công bố của chính lớp đang xếp

Mục tiêu:

- nếu xếp lại chính lớp đó và gặp buổi `PUBLISHED` đã có từ trước thì không báo trùng toàn bộ như conflict ngoài;
- lịch mới có thể đè lên chính buổi cũ theo cơ chế replan;
- vẫn báo conflict nếu trùng giáo viên/phòng/lớp với buổi published của lớp khác hoặc đối tượng khác.

Task:

- trong conflict detector, phân biệt `same_class_same_session/replan_target` với conflict thật;
- nếu existing lesson thuộc class đang replan, map thành candidate có thể thay thế;
- thêm action metadata: `replace_lesson_id`;
- commit preview cập nhật lesson cũ hoặc tạo draft thay thế theo policy đã chọn;
- thêm test regression cho class tự đè lesson published của chính nó.

Acceptance:

- replan một lớp đã có lịch published không sinh conflict hàng loạt với chính lịch của nó;
- vẫn chặn trùng phòng/giáo viên với lịch của lớp khác;
- UI hiển thị đây là "cập nhật buổi đã có", không phải "trùng lịch".

### Nhóm P1 - Đơn giản hóa UI preview

P1.1. Bỏ bộ lọc chế độ preview trên UI

Mục tiêu:

- UI vận hành không còn chọn `cold_start/replan_draft/replan_with_published_lock`;
- mặc định gửi mode giữ lịch đã công bố.

Task:
- field chọn theo lớp/giảng viên chuyển sang thành radio button;
- bỏ field chọn chế độ khỏi form;
- hardcode request mode là giữ lịch đã công bố;
- giữ mode khác cho benchmark/dev tool nếu cần;
- bỏ subtext giải thích mode khỏi màn preview.

Acceptance:

- người dùng chỉ chọn đối tượng và ngày bắt đầu dự kiến;
- request preview vẫn gửi đúng mode mặc định;
- không còn text "cold start", "draft", "published lock" trên UI.

P1.2. Dọn text, subtext, alert rườm rà

Mục tiêu:

- màn preview gọn để giáo vụ thao tác nhanh;
- không còn nửa Anh nửa Việt.

Task:

- rà toàn bộ `SchedulingPage.tsx` bằng search các từ tiếng Anh đang hiển thị;
- bỏ tổng quan `run_id` preview ở header;
- bỏ ô tìm lớp/giáo viên trong preview lịch nếu không dùng thật;
- bỏ alert chỉ mô tả hiển nhiên;
- chỉ giữ alert lỗi thật, cảnh báo trùng thật, cảnh báo không thể lưu;
- chuyển các giải thích dài sang tooltip.

Acceptance:

- UI preview chỉ còn: cấu hình ngắn, summary từng lớp, calendar, danh sách ca cần xử lý;
- không còn text tiếng Anh trên màn;
- số alert trong trạng thái bình thường gần như bằng 0.

P1.3. Calendar theo khung giờ ca

Mục tiêu:

- lịch tuần nhìn như bảng ca theo hàng ngang;
- ô trống để trống, không ghi "không có buổi".

Task:

- thêm trục dọc là ca/khung giờ: ví dụ `07:00-09:00`, `18:00-20:00`;
- cột ngang là thứ/ngày;
- card buổi học nằm đúng hàng ca và cột ngày;
- ô trống không render text;
- chip/card không cần lặp lại ngày giờ nếu hàng/cột đã thể hiện;
- tooltip card hiển thị ngày giờ đầy đủ, buổi thứ mấy, phòng, giáo viên.

Acceptance:

- nhìn ngang biết ngay một ca có những lớp nào;
- ô trống không còn text dư;
- card gọn hơn nhưng hover vẫn đủ thông tin.

P1.4. Chỉnh tay theo luồng chọn ca rồi chọn phòng

Mục tiêu:

- giảm danh sách phương án quá dài;
- người dùng chọn theo tư duy vận hành thực tế.

Luồng mới:

1. Chọn ca.
2. Chọn phòng.
3. Hệ thống hiển thị các ngày/vị trí hợp lệ còn lại.
4. Chọn slot.

Task:

- group `candidate_options` theo `shift_id`;
- sau khi chọn ca, lọc tiếp phòng hợp lệ trong ca đó;
- sau khi chọn phòng, chỉ hiện các ngày/slot hợp lệ;
- tooltip slot hiển thị điểm mềm, điểm cộng/trừ và buổi thứ mấy;
- test lại kéo thả với dataTransfer và state giữ đúng lúc drop.

Acceptance:

- dialog đổi ca không còn danh sách tất cả ca/phòng trộn chung;
- số phương án nhìn thấy giảm mạnh;
- kéo thả áp dụng được phương án đã chọn;
- điểm mềm có giải thích đủ hiểu.

P1.5. Summary từng lớp thay cho tổng quan preview run

Mục tiêu:

- bỏ tổng quan `id preview`;
- trạng thái xem trước nằm trên từng lớp.

Task:

- mỗi lớp hiển thị: trạng thái, số buổi đã xếp/tổng buổi, số trùng, số ca cần xử lý;
- nếu có trùng, thêm chip ngay sau nhãn lớp;
- bỏ panel header preview run nếu không cần debug.

Acceptance:

- người dùng nhìn từng lớp biết lớp nào cần xử lý;
- không cần đọc `run_id`;
- thông tin tổng chỉ giữ ở dạng rất gọn nếu cần.

### Nhóm P2 - Chuẩn hóa kỹ năng/chứng chỉ

P2.1. Tách kỹ năng/chứng chỉ thành bảng riêng

Mục tiêu:

- kỹ năng là dữ liệu chuẩn hóa, có thể quản trị, lọc và tái sử dụng;
- dùng chung cho phân công giáo viên, preview và dạy thay.

Model đề xuất:

- `skills`: `id`, `code`, `name`, `description`, `is_active`;
- `teacher_skills`: `teacher_id`, `skill_id`, `level`, `verified_at`, `expires_at`;
- `course_required_skills`: `course_id`, `skill_id`, `required_level`, `is_mandatory`;

Task:

- tạo migration mới;
- migrate dữ liệu từ field mảng hiện tại sang bảng mới;
- giữ compatibility tạm thời nếu cần;
- cập nhật repository/usecase/form admin.

Acceptance:

- tạo/sửa skill được;
- gắn skill cho giáo viên được;
- gắn skill yêu cầu cho khóa học được;
- preview không đọc field mảng cũ nữa.

P2.2. Validate khi gán giáo viên vào lớp

Mục tiêu:

- không gán giáo viên sai kỹ năng cho lớp;
- danh sách giáo viên chọn trong lớp được lọc sẵn.

Task:

- API list giáo viên nhận `course_id` hoặc `class_id` để lọc giáo viên đủ kỹ năng;
- use case assign teacher validate skill trước khi lưu;
- UI class detail chỉ hiển thị giáo viên phù hợp, hoặc hiển thị nhóm "không phù hợp" nếu cần debug;
- message lỗi thuần tiếng Việt.

Acceptance:

- không lưu được giáo viên thiếu skill bắt buộc;
- danh sách chọn giáo viên gọn hơn;
- lỗi rõ thiếu kỹ năng nào.

P2.3. Dùng skill trong đề xuất dạy thay

Mục tiêu:

- dạy thay chỉ đề xuất giáo viên cùng kỹ năng/chứng chỉ.

Task:

- substitute matcher đọc `course_required_skills`;
- loại giáo viên thiếu skill bắt buộc;
- score thêm điểm cho giáo viên skill cao hơn hoặc chứng chỉ còn hiệu lực lâu;
- tooltip lý do đề xuất hiển thị skill match.

Acceptance:

- danh sách dạy thay không có giáo viên sai chuyên môn;
- mỗi gợi ý có lý do: kỹ năng, lịch trống, tải công việc, di chuyển.

### Nhóm P3 - Quản lý buổi học

P3.1. Bỏ filter lifecycle trên UI quản lý buổi học

Mục tiêu:

- lifecycle là chi tiết kỹ thuật, không phải bộ lọc chính của giáo vụ.

Task:

- bỏ filter lifecycle khỏi `LessonsPage`;
- giữ filter ngày, lớp, giáo viên, trạng thái nghiệp vụ nếu có;
- giản lược subtitle/subtext;
- nếu cần hiện lifecycle, chỉ hiện chip nhỏ trên dòng dữ liệu hoặc tooltip.

Acceptance:

- màn buổi học gọn hơn;
- không còn filter "cyclelife/lifecycle";
- text thuần tiếng Việt.

## 4. Thứ tự triển khai khuyến nghị mới

### Đợt 5A - Sửa đúng sai scheduling core

Ưu tiên làm trước UI vì các lỗi này ảnh hưởng kết quả xếp lịch.

Task:

- P0.1. Không xếp lịch về quá khứ;
- P0.2. Window và số buổi riêng từng lớp;
- P0.3. Không báo trùng với published lesson của chính lớp đang replan;
- thêm test regression cho cả 3 case.

Done khi:

- preview không sinh lịch quá khứ;
- chọn nhiều lớp vẫn đúng số buổi từng lớp;
- replan lớp đã có lịch không tự conflict với chính nó.

### Đợt 5B - Gọt UI preview thành bản vận hành gọn

Task:

- P1.1. Bỏ bộ lọc chế độ preview;
- P1.2. Dọn text, subtext, alert;
- P1.5. Summary từng lớp thay tổng quan run;
- rà toàn bộ text tiếng Anh trong `SchedulingPage`.

Done khi:

- màn preview thuần tiếng Việt;
- ít alert, ít subtitle;
- không còn `run_id` làm thông tin chính;
- trạng thái lớp và số trùng nằm ngay trên summary lớp.

### Đợt 5C - Làm lại calendar và chỉnh tay

Task:

- P1.3. Calendar theo khung giờ ca;
- P1.4. Chọn ca -> chọn phòng -> chọn slot;
- tooltip buổi thứ mấy, điểm cộng/trừ;
- test kéo thả thật trên browser.

Done khi:

- calendar có hàng khung giờ;
- ô trống không có text;
- chọn slot thủ công ít rối hơn;
- kéo thả dùng được.

### Đợt 5D - Chuẩn hóa skill/chứng chỉ

Task:

- P2.1. Tạo bảng skill;
- P2.2. Validate và lọc giáo viên khi gán lớp;
- P2.3. Dùng skill chuẩn hóa trong dạy thay.

Done khi:

- skill không còn là mảng text tự do;
- assign teacher chặn sai kỹ năng;
- substitute suggestion dựa trên skill chuẩn.

### Đợt 5E - Dọn quản lý buổi học

Task:

- P3.1. Bỏ lifecycle filter;
- giản lược subtitle/subtext;
- rà text tiếng Việt.

Done khi:

- màn quản lý buổi học gọn và bớt thuật ngữ kỹ thuật.

## 5. Test và benchmark cần bổ sung

Backend tests:

- preview với ngày bắt đầu trong quá khứ;
- nhiều lớp có số buổi khác nhau;
- giáo viên dạy nhiều lớp có window riêng từng lớp;
- replan chính lớp có published lesson cũ;
- assign teacher thiếu skill bị chặn;
- substitute không trả giáo viên thiếu skill.

Frontend tests hoặc manual QA:

- màn preview không còn text tiếng Anh;
- bỏ mode filter vẫn gửi đúng mode mặc định;
- calendar có hàng khung giờ;
- ô trống không render text;
- dialog đổi slot theo đúng luồng ca -> phòng -> slot;
- tooltip hiển thị buổi thứ mấy, điểm cộng, điểm trừ;
- kéo thả áp dụng được slot.

Benchmark:

- cold start chỉ dùng để benchmark nội bộ;
- production preview mặc định giữ lịch đã công bố;
- báo cáo cần có: thời gian xử lý, số conflict luật cứng, số ca cần chỉnh tay, điểm mềm, số ca đổi lịch, mức lấp đầy.

## 6. Definition of Done cho mốc kế tiếp

Mốc kế tiếp đạt yêu cầu khi:

- không còn sinh lịch quá khứ;
- số buổi từng lớp đúng theo khóa học;
- preview không tự báo trùng với lịch published của chính lớp đang xếp lại;
- màn preview đủ gọn để giáo vụ dùng hằng ngày;
- toàn bộ text UI liên quan scheduling là tiếng Việt;
- chỉnh tay theo ca, phòng, slot;
- skill/chứng chỉ có bảng chuẩn hóa và được dùng khi gán giáo viên, xếp lịch, dạy thay.
