### Mục tiêu
Khởi tạo Git cho dự án hiện tại và đẩy lên GitHub bằng tài khoản cá nhân, trong bối cảnh bạn dùng nhiều profile Git (nhiều tài khoản/SSH key).

Dưới đây là quy trình an toàn, không ảnh hưởng tới các repo khác. Bạn chỉ cần làm một lần cho repo này.

---

### 1) Vào thư mục dự án
```bash
cd "/Users/hant/golang/doan"
```

### 2) Khởi tạo Git repo và tạo nhánh chính `main`
```bash
git init
# Nếu Git cũ tạo mặc định master, đổi sang main
git branch -M main
```

### 3) Thiết lập profile Git CỤC BỘ cho repo này (tài khoản cá nhân)
Thiết lập `user.name` và `user.email` chỉ trong repo này (không ảnh hưởng máy toàn cục):
```bash
git config user.name "Tên Cá Nhân"
git config user.email "email_ca_nhan@domain.com"
# Tuỳ chọn: buộc Git chỉ dùng config cục bộ thay vì rơi về global
git config user.useConfigOnly true
```

Nếu bạn ký commit bằng GPG/SSH, mình có thể hướng dẫn thêm sau khi bạn xác nhận bạn dùng loại nào. Với hầu hết trường hợp, 3 dòng trên là đủ.

---

### 4) Tạo (hoặc dùng) SSH key riêng cho tài khoản GitHub cá nhân
Nếu bạn đã có SSH key cá nhân, chuyển sang bước 5. Nếu chưa, tạo mới:
```bash
ssh-keygen -t ed25519 -C "email_ca_nhan@domain.com" -f ~/.ssh/id_ed25519_personal
# bật ssh-agent và thêm key
eval "$(ssh-agent -s)"
ssh-add ~/.ssh/id_ed25519_personal
```

Tạo cấu hình SSH để tách bạch nhiều tài khoản (rất quan trọng khi dùng nhiều profile):
```bash
mkdir -p ~/.ssh
chmod 700 ~/.ssh
cat >> ~/.ssh/config <<'EOF'
Host github.com-personal
  HostName github.com
  User git
  IdentityFile ~/.ssh/id_ed25519_personal
  IdentitiesOnly yes
EOF
chmod 600 ~/.ssh/config
```

Thêm public key lên GitHub cá nhân:
- Lấy public key:
  ```bash
  pbcopy < ~/.ssh/id_ed25519_personal.pub  # macOS (hoặc: cat ~/.ssh/id_ed25519_personal.pub)
  ```
- Vào GitHub → Settings → SSH and GPG keys → New SSH key → dán nội dung public key → Save.

Kiểm tra kết nối với alias cá nhân:
```bash
ssh -T git@github.com-personal
# Kỳ vọng: "Hi <username>! You've successfully authenticated, but GitHub does not provide shell access."
```

---

### 5) Tạo repo trống trên GitHub
Bạn có thể tạo thủ công trên web: New → Repository → Name: `doan` (hoặc tên bạn muốn) → Private/Public tuỳ nhu cầu → KHÔNG tạo README/.gitignore (vì đã có local). Ghi nhớ `username` của bạn.

Hoặc dùng GitHub CLI (nếu đã cài `gh` và login bằng tài khoản cá nhân):
```bash
# Tạo repo từ thư mục hiện tại, gắn remote và push ngay
# Tip: gh sẽ đặt remote https/ssh theo host mặc định của tài khoản đã login trong gh
# Nếu bạn dùng nhiều tài khoản, cách chắc chắn nhất vẫn là tự đặt remote (bước 6) thay vì rely vào gh.
```

---

### 6) Thêm remote dùng đúng profile cá nhân (SSH alias)
Thay `your-github-username` bằng username GitHub cá nhân của bạn.
```bash
git remote add origin git@github.com-personal:your-github-username/doan.git
```

Nếu lỡ đã có `origin`, cập nhật lại:
```bash
git remote set-url origin git@github.com-personal:your-github-username/doan.git
```

Kiểm tra:
```bash
git remote -v
```

---

### 7) Thêm .gitignore phù hợp (Go)
Nếu bạn chưa có `.gitignore`, tạo nhanh theo mẫu gợi ý cho dự án Go (kèm IDE phổ biến):
```gitignore
# Go build
bin/
build/
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
coverage*

# Go modules
vendor/

# Logs & temp
*.log
.tmp/
.cache/

# IDE/editor
.idea/
.vscode/
*.swp

# OS
.DS_Store
```
Bạn có thể điều chỉnh theo dự án (ví dụ nếu bạn cần `vendor` thì bỏ dòng `vendor/`).

Thêm và commit:
```bash
git add .
git commit -m "chore: initial commit"
```

---

### 8) Đẩy lên GitHub
```bash
git push -u origin main
```
Sau đó mở URL repo trên GitHub để kiểm tra.

---

### 9) Kiểm tra đúng profile đang dùng
- Xem user/email của repo hiện tại:
  ```bash
  git config user.name
  git config user.email
  ```
- Xác nhận remote dùng alias cá nhân:
  ```bash
  git remote -v
  # phải thấy: git@github.com-personal:your-github-username/doan.git
  ```
- Test SSH đúng key:
  ```bash
  ssh -vT git@github.com-personal  # -v hiển thị debug key nào đang dùng
  ```

---

### Tuỳ chọn: Tổ chức nhiều profile theo thư mục (nâng cao)
Nếu bạn muốn tự động hoá user/email theo thư mục (vd: tất cả repo trong `~/personal/` dùng profile cá nhân):

1) Tạo hai file cấu hình:
```ini
# ~/.gitconfig-personal
[user]
  name = Tên Cá Nhân
  email = email_ca_nhan@domain.com
```
```ini
# ~/.gitconfig-work
[user]
  name = Tên Công Ty
  email = email_cong_ty@domain.com
```
2) Thêm include có điều kiện vào `~/.gitconfig`:
```ini
[includeIf "gitdir:~/personal/"]
  path = ~/.gitconfig-personal
[includeIf "gitdir:~/work/"]
  path = ~/.gitconfig-work
```
3) Đặt repo vào thư mục tương ứng (`~/personal/doan`), Git sẽ tự áp dụng profile. SSH alias vẫn dùng như mục 4.

---

### Lỗi thường gặp & cách xử lý nhanh
- `Permission denied (publickey)`: Chưa thêm public key lên GitHub hoặc remote không dùng đúng alias. Sửa remote sang `git@github.com-personal:...` và kiểm tra `ssh -T`.
- Push lên nhầm tài khoản: Thường do dùng `git@github.com:...` (mặc định) thay vì alias `github.com-personal`. Đổi remote theo mục 6.
- Không đúng `user.email` trên commit: Kiểm tra `git config user.email`. Dùng `git commit --amend --reset-author` để cập nhật commit mới nhất nếu cần.
- Đã có lịch sử trước khi tạo repo GitHub: Cứ add remote và `git push -u origin main` bình thường.

---

### Bạn cần mình thực hiện hộ bước nào không?
Hãy cho mình:
- Username GitHub cá nhân bạn muốn dùng
- Email cá nhân dùng cho commit
- Bạn đã có SSH key cá nhân chưa (và đường dẫn), hay muốn tạo mới theo hướng dẫn trên?

Mình sẽ trả lời lại với các lệnh đã điền sẵn giá trị chính xác cho bạn chạy copy-paste.