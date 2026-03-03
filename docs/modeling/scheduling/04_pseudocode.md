# Backtracking CSP Scheduling - Pseudocode

Mục tiêu của tài liệu này là cung cấp giả mã (Pseudocode) chi tiết cho thuật toán cốt lõi điều phối thời khóa biểu.

---

## 1. Hàm chính: `SolveCSP`
Hàm khởi tạo và kích hoạt Backtracking:

```pascal
function SolveCSP(Variables, Domains):
    Initialize empty Assignment = {}
    return Backtrack(Assignment, Variables, Domains)
```

## 2. Hàm Đệ quy: `Backtrack`
Đây là lõi thuật toán tìm kiếm theo chiều sâu (DFS search) qua cây trạng thái gán.

```pascal
function Backtrack(Assignment, Variables, Domains):
    // Cơ sở đệ quy: Nếu mọi lớp đều đã được gán lịch
    if isValidAndComplete(Assignment) == true:
        return Assignment

    // Heuristic MRV: Chọn biến (Class/Session) có số lượng Domain (Option) ít nhất
    variable = SelectUnassignedVariable(Variables, Assignment, Domains)
    
    // Iteration qua các Domain Values
    // Heuristic LCV: Ưu tiên thử các giá trị ít gây xung đột nhất cho các biến khác
    for each value in OrderDomainValues(variable, Domains, Assignment):
        
        // Hard Constraint Check (Không trùng phòng, không trùng giáo viên, <= 22h, v.v)
        if CheckHardConstraints(variable, value, Assignment) == true:
            
            // Tạm thời gán giá trị
            Add {variable = value} to Assignment
            
            // Tạm thời thu hẹp (Pruning) miền giá trị của các biến liên quan để loại bỏ xung đột (Forward Checking)
            inferences = ForwardCheck(variable, value, Variables, Domains, Assignment)
            
            // Nếu việc thu hẹp không làm bất kỳ biến nào bị rỗng Domain (nghĩa là không dẫn tới Failure vô vọng)
            if inferences != FAILURE:
                
                // Cập nhật Domain và tiếp tục đi sâu xuống nhánh này
                ApplyInferences(Domains, inferences)
                result = Backtrack(Assignment, Variables, Domains)
                
                if result != FAILURE:
                    return result
                
                // Undo Domain nếu nhánh phân giải thất bại
                RemoveInferences(Domains, inferences)
            
            // Undo gán giá trị (Quay lui - Backtrack)
            Remove {variable = value} from Assignment

    return FAILURE
```

## 3. Hàm Heuristics: `SelectUnassignedVariable` (MRV)
Mininum Remaining Values (MRV).

```pascal
function SelectUnassignedVariable(Variables, Assignment, Domains):
    min_domain_size = INFINITY
    best_variable = NULL
    
    for each v in Variables:
        if v not in Assignment:
            current_domain_size = countValidValues(v, Domains[v], Assignment)
            if current_domain_size < min_domain_size:
                min_domain_size = current_domain_size
                best_variable = v
            // Bẻ khóa (tie breaker) bằng Degree Heuristic nếu cần
            else if current_domain_size == min_domain_size:
                best_variable = DegreeHeuristic(v, best_variable)
                
    return best_variable
```

## 4. Hàm Forward Checking: `ForwardCheck`
Kiểm tra xem khi gán `v = d`, liệu có Variable láng giềng nào mất hoàn toàn Options không.

```pascal
function ForwardCheck(variable, assigned_value, Variables, Domains, Assignment):
    inferences = Initialize Empty Map
    
    // Tìm các Session có khả năng xung đột (Ví dụ: cùng Room, cùng Teacher)
    neighbors = GetNeighbors(variable) 
    
    for each neighbor in neighbors:
        if neighbor not in Assignment:
            valid_values_for_neighbor = []
            
            for each value in Domains[neighbor]:
                // Giả lập đưa vào và kiểm tra
                if CheckHardConstraints(neighbor, value, Assignment U {variable=assigned_value}):
                    valid_values_for_neighbor.add(value)
            
            if valid_values_for_neighbor is EMPTY:
                // Một biến láng giềng không còn lựa chọn nào -> Nhánh này vô hiệu hoá
                return FAILURE
            
            // Ghi nhận miền giá trị được thu hẹp
            inferences[neighbor] = valid_values_for_neighbor
            
    return inferences
```

## 5. Hàm Tính điểm Ràng buộc Mềm (Soft Constraint Scoring)
Hàm này chạy ngầm phía dưới khi `OrderDomainValues` gọi đến để ưu tiên gán LCV.

```pascal
function ScoreSoftConstraints(Assignment):
    score = 0
    
    // Phân tích giáo viên
    for each teacher in AllTeachers():
        sessions = GetSessionsForTeacher(Assignment, teacher)
        sessions = SortByTime(sessions)
        
        for i = 0 to length(sessions) - 2:
            gap = sessions[i+1].startTime - sessions[i].endTime
            
            if gap <= 30 minutes:
                score += 10 // Thưởng (Liên tiếp)
            else if gap > 120 minutes:
                score -= round(gap / 60) * 5 // Phạt (Trống lịch quá dài)
                
    return score
```
