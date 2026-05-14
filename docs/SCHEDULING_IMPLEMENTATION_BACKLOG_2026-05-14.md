# Scheduling Implementation Backlog 2026-05-14

Tài liệu này bóc phần định hướng scheduling production thành backlog triển khai bám sát codebase hiện tại.

## 1. Trạng thái hiện tại trong repo

Những phần đã có:

- preview scheduling và commit preview sang `lesson`;
- benchmark so sánh solver, trong đó `CP-SAT` đang là solver chính;
- rule lọc lớp chưa đủ sĩ số tối thiểu trước khi đưa vào preview;
- rule chặn commit nếu số ca chỉnh tay vượt ngưỡng vận hành;
- UI preview có kéo thả và hiển thị `score delta` cho lựa chọn chỉnh tay;
- màn lịch admin, teacher, student đã bắt đầu dùng cùng ngôn ngữ hiển thị theo tuần.

Những phần chưa có trong model/code:

- chưa có trạng thái vòng đời `History/Published/Draft/Unplanned` trên `lesson`;
- chưa có mô hình `campus`, thời gian di chuyển, kỹ năng giáo viên và kỹ năng yêu cầu của khóa học;
- chưa có use case production cho dạy thay, học bù, chuyển lớp, tái lập lịch ít xáo trộn;
- chưa có benchmark production-like cho substitution, makeup và replanning;
- route UI `teacher/substitute` vẫn đang là placeholder.

## 2. Kết luận nghiên cứu ngắn

Nếu triển khai thẳng các luật mới vào solver ngay bây giờ thì sẽ bị thiếu dữ liệu đầu vào.

Vì vậy cần đi theo thứ tự:

1. mở rộng domain và migration;
2. chuẩn hóa API và use case vận hành;
3. cài ràng buộc mới vào solver;
4. hoàn thiện UI;
5. thêm benchmark và test.

## 3. Epic A - Chuẩn hóa domain và dữ liệu nền

### A1. Bổ sung vòng đời lesson

Mục tiêu:

- cho phép phân biệt `History`, `Published`, `Draft`, `Unplanned`;
- làm nền cho continuous planning và nonvolatile replanning.

Task:

- thêm các field trạng thái vào `Lesson`;
- thêm metadata như `published_at`, `draft_group_id`, `source_preview_run_id`, `change_reason`;
- thống nhất rule chuyển trạng thái;
- cập nhật DTO, repository, migration và seed.

Acceptance:

- lesson mới từ preview commit có trạng thái rõ ràng;
- API list lesson lọc được theo trạng thái;
- `History` không bị đưa vào solver.

Dependency:

- không có, nên làm đầu tiên.

### A2. Bổ sung mô hình campus và travel time

Mục tiêu:

- hỗ trợ luật di chuyển giữa các cơ sở.

Task:

- thêm `Campus` hoặc tối thiểu `campus_id` cho `Room`;
- thêm bảng `campus_travel_times` hoặc cấu hình matrix thời gian di chuyển;
- cập nhật CRUD room và seed dữ liệu;
- thêm helper tính travel feasibility giữa hai lesson.

Acceptance:

- mỗi room có thể truy ra campus;
- hệ thống tính được thời gian di chuyển tối thiểu giữa 2 campus.

Dependency:

- độc lập với lesson lifecycle, nhưng nên làm trước khi viết constraint travel.

### A3. Bổ sung kỹ năng giáo viên và yêu cầu kỹ năng khóa học

Mục tiêu:

- hỗ trợ hard constraint skill matching;
- làm nền cho substitute matching.

Task:

- thêm `TeacherSkill` hoặc field chứng chỉ cho giáo viên;
- thêm `required_skills` cho `Course`;
- chuẩn hóa mapper giữa kỹ năng của teacher và yêu cầu của course;
- cập nhật form admin và seed dữ liệu benchmark.

Acceptance:

- hệ thống xác định được teacher có đủ kỹ năng cho class hay không;
- preview conflict hiện được lỗi thiếu kỹ năng nếu teacher không hợp lệ.

Dependency:

- nên làm trước substitute use case và solver upgrade.

## 4. Epic B - Nâng cấp scheduling core

### B1. Tách solver theo vùng Draft và Published

Mục tiêu:

- production scheduling chỉ chạy trong vùng `Draft`;
- `Published` là semi-movable entity.

Task:

- sửa input builder của preview để nạp lesson theo lifecycle;
- xem `Published` như lesson existing đặc biệt;
- loại `History` khỏi vùng tái lập lịch;
- thêm option preview mode: `cold_start`, `replan_draft`, `replan_with_published_lock`.

Acceptance:

- preview không động vào lesson `History`;
- preview có thể khóa cứng lesson `Published`;
- commit conflict message phân biệt rõ conflict với published lesson.

Dependency:

- cần A1.

### B2. Cài hard constraint travel time

Mục tiêu:

- không xếp giáo viên di chuyển không kịp giữa 2 campus.

Task:

- mở rộng conflict detector và candidate filter với travel feasibility;
- thêm hard check trong solver input/domain validation;
- thêm conflict type mới, ví dụ `TRAVEL_TIME_BLOCK`.

Acceptance:

- candidate option vi phạm travel time không được chọn;
- preview hiển thị conflict travel rõ ràng.

Dependency:

- cần A2.

### B3. Cài hard constraint skill matching

Mục tiêu:

- chỉ teacher có kỹ năng phù hợp mới được dạy class hoặc nhận dạy thay.

Task:

- validate teacher of class khi tạo variable;
- chặn candidate/substitute nếu teacher không đủ skill;
- thêm conflict type `SKILL_MISMATCH`.

Acceptance:

- class dùng teacher sai skill không đi qua preview clean;
- substitute suggestion không trả teacher sai chuyên môn.

Dependency:

- cần A3.

### B4. Nonvolatile replanning

Mục tiêu:

- khi replan, hệ thống ưu tiên ít xáo trộn lesson đã công bố.

Task:

- thêm penalty rất lớn khi đổi slot/phòng/teacher của lesson `Published`;
- thêm metric `schedule_change_count`, `teacher_change_count`, `room_change_count`;
- hiển thị impact summary trước commit.

Acceptance:

- cùng một input sự cố, solver ưu tiên sửa cục bộ thay vì xáo toàn lịch;
- preview trả được số lượng lesson bị ảnh hưởng.

Dependency:

- cần A1 và B1.

### B5. Capacity utilization và rule học bù/chuyển lớp

Mục tiêu:

- dùng sức chứa theo thời gian thực cho makeup và transfer.

Task:

- tính `capacity_utilization` cho lesson/lớp/phòng;
- bổ sung helper tìm slot còn chỗ;
- chặn makeup/transfer nếu vượt sức chứa tại thời điểm lesson diễn ra.

Acceptance:

- API tìm học bù trả được số chỗ còn lại;
- không thể xác nhận nhét thêm học viên nếu vượt capacity.

Dependency:

- dùng được với model hiện tại, nhưng tốt hơn nếu hoàn thành A1 trước.

## 5. Epic C - Use case vận hành thời gian thực

### C1. Dạy thay trong 30 giây

Mục tiêu:

- khi teacher nghỉ đột xuất, hệ thống đề xuất người dạy thay khả dụng.

Task:

- tạo use case `SuggestSubstituteTeachers`;
- input gồm lesson, required skill, thời gian, campus;
- rank theo skill match, lịch trống, workload balance, travel feasibility;
- thêm API admin/teacher cho request và duyệt dạy thay.

Acceptance:

- từ một lesson có thể gọi API trả danh sách teacher thay thế;
- mỗi teacher có score và lý do gợi ý;
- thời gian xử lý đạt ngưỡng vận hành mục tiêu.

Dependency:

- cần A2, A3, B2, B3.

### C2. Find spot cho học bù

Mục tiêu:

- học viên nghỉ có thể được nhét vào lesson tương thích còn chỗ.

Task:

- tạo use case `FindMakeupSpots`;
- lọc lớp cùng course hoặc cùng level;
- kiểm tra student conflict, capacity, campus, thời gian;
- hỗ trợ xác nhận học bù vào lesson đích.

Acceptance:

- API trả được danh sách lesson học bù hợp lệ;
- có thể tạo record học bù mà không vi phạm sức chứa hoặc trùng lịch student.

Dependency:

- cần B5.

### C3. Bảo lưu và chuyển lớp

Mục tiêu:

- thay đổi roster không phá vỡ lịch đang chạy.

Task:

- cập nhật rule khi học viên bảo lưu;
- khi chuyển lớp, kiểm tra class đích còn chỗ và cùng điều kiện;
- nếu cần, đánh dấu nhu cầu `Unplanned` cho các buổi phải học bù.

Acceptance:

- luồng chuyển lớp không làm dữ liệu scheduling bị âm thầm sai;
- số chỗ trống và utilization cập nhật đúng sau nghiệp vụ.

Dependency:

- cần A1 và B5.

## 6. Epic D - API và UI

### D1. Mở rộng scheduling API

Task:

- thêm field lifecycle và impact summary vào preview response;
- thêm endpoint replanning;
- thêm endpoint substitute suggestion;
- thêm endpoint find makeup spot;
- thêm endpoint publish draft lessons.

Acceptance:

- frontend không phải tự suy diễn trạng thái lesson nữa;
- API đủ để dựng flow production.

### D2. Hoàn thiện màn admin scheduling

Task:

- thêm bộ lọc `Draft/Published/History/Unplanned`;
- hiển thị change impact khi replan;
- hiển thị cảnh báo travel/skill/capacity;
- thêm panel summary cho `schedule_change_count` và `capacity_utilization`.

Acceptance:

- người xếp lịch thấy ngay tác động của một lần replan;
- preview không chỉ báo hard conflict mà còn báo ảnh hưởng vận hành.

### D3. Thay placeholder bằng màn hình dạy thay

Task:

- thay route `teacher/substitute` placeholder bằng page thật;
- teacher có thể gửi request nghỉ/dạy thay;
- admin có thể duyệt teacher thay thế.

Acceptance:

- route dạy thay dùng được end-to-end.

### D4. Màn hình học bù/chuyển lớp

Task:

- thêm modal hoặc page để tìm slot học bù;
- cho phép xác nhận nhét học viên vào lesson tương thích;
- hiển thị remaining capacity và cảnh báo trùng lịch.

Acceptance:

- giáo vụ có thể tìm và chốt slot học bù mà không cần rà tay lịch.

## 7. Epic E - Benchmark, seed và test

### E1. Seed large dataset chuẩn

Task:

- tạo dataset lớn có nhiều lớp, teacher, student, room, campus;
- có teacher skill và course required skill;
- có published draft history mix;
- có scenario teacher nghỉ, room hỏng, học viên học bù.

Acceptance:

- benchmark và demo dùng cùng một bộ seed có thể tái tạo.

### E2. Benchmark production-like

Task:

- thêm scenario `substitution`;
- thêm scenario `makeup_find_spot`;
- thêm scenario `published_stability`;
- thêm metric `schedule_change_count`, `capacity_utilization`, `score_delta`.

Acceptance:

- có report số liệu để bảo vệ quyết định chọn `CP-SAT` cho bài toán production-like, không chỉ cold start.

### E3. Test coverage

Task:

- unit test cho policy mới;
- integration test cho preview/replan/substitute/makeup;
- UI test cho kéo thả, delta score, lifecycle filter.

Acceptance:

- mỗi epic đều có test regression tương ứng;
- benchmark artifacts chạy được trong repo.

## 8. Thứ tự triển khai khuyến nghị

### Đợt 1 - Chốt nền dữ liệu

- A1. Lesson lifecycle
- A2. Campus và travel time
- A3. Teacher skill và course required skill

### Đợt 2 - Chốt scheduling core

- B1. Draft/Published scheduling boundary
- B2. Travel hard constraint
- B3. Skill matching hard constraint
- B4. Nonvolatile replanning
- B5. Capacity utilization

### Đợt 3 - Làm production operations

- C1. Dạy thay
- C2. Học bù
- C3. Chuyển lớp và bảo lưu

### Đợt 4 - Hoàn thiện API/UI

- D1. Scheduling API
- D2. Admin scheduling UX
- D3. Teacher substitute page
- D4. Makeup/transfer UI

### Đợt 5 - Benchmark và chốt số liệu

- E1. Seed chuẩn
- E2. Benchmark production-like
- E3. Test coverage

## 9. Lát cắt nên làm tiếp ngay

Nếu cần chọn nhánh triển khai tiếp theo ngay sau tài liệu này, nên làm:

1. `A1` vì nó mở khóa replanning và publish stability.
2. `A3` vì skill matching ảnh hưởng cả preview lẫn substitute.
3. `D1` phần lifecycle trong preview response để UI và backend cùng nói chung một ngôn ngữ.

## 10. Định nghĩa Done cho phiên bản production-like đầu tiên

Có thể xem scheduling đạt mốc production-like đầu tiên khi thỏa đủ:

- preview chỉ chạy trên vùng `Draft`;
- lesson có lifecycle rõ ràng;
- có rule mở lớp theo sĩ số tối thiểu;
- có travel time và skill matching trong hard constraints;
- có `score delta` và impact summary khi chỉnh tay hoặc replan;
- có flow dạy thay và tìm học bù cơ bản;
- có benchmark production-like dùng số liệu thật để giải thích vì sao tiếp tục chọn `CP-SAT`.
