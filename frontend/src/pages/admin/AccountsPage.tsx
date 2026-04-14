import { useMemo, useState } from 'react';
import {
    Alert,
    Box,
    Button,
    Chip,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    FormControl,
    InputLabel,
    MenuItem,
    Paper,
    Select,
    Stack,
    Switch,
    TextField,
    Typography,
} from '@mui/material';
import { Add, Edit, LockReset, SearchRounded, Visibility } from '@mui/icons-material';
import { DataGrid } from '@mui/x-data-grid';
import type { GridColDef, GridRenderCellParams } from '@mui/x-data-grid';
import { toast } from 'sonner';

import {
    useCreateUserMutation,
    useGetUsersQuery,
    useLazyGetUserByIdQuery,
    useResetUserPasswordMutation,
    useUpdateUserMutation,
    type NormalizedAccountUser,
} from '@/api/accountApi';
import PageHeader from '@/components/common/PageHeader';
import ConfirmDialog from '@/components/common/ConfirmDialog';
import { getApiErrorMessage } from '@/utils/apiError';

const roles = ['SUPER_ADMIN', 'ADMIN', 'TEACHER', 'STUDENT', 'PARENT', 'COMPLIANCE'];

type AccountFormState = {
    code: string;
    full_name: string;
    email: string;
    role: string;
    is_active: boolean;
    password: string;
};

const emptyForm: AccountFormState = {
    code: '',
    full_name: '',
    email: '',
    role: 'STUDENT',
    is_active: true,
    password: '',
};

export const AccountsPage = () => {
    const [page, setPage] = useState(0);
    const [pageSize, setPageSize] = useState(10);
    const [searchInput, setSearchInput] = useState('');
    const [search, setSearch] = useState('');
    const [role, setRole] = useState('');
    const [onlyActive, setOnlyActive] = useState(false);

    const [dialogOpen, setDialogOpen] = useState(false);
    const [detailOpen, setDetailOpen] = useState(false);
    const [resetDialogOpen, setResetDialogOpen] = useState(false);
    const [selectedUser, setSelectedUser] = useState<NormalizedAccountUser | null>(null);
    const [form, setForm] = useState<AccountFormState>(emptyForm);
    const [resetPasswordValue, setResetPasswordValue] = useState('');

    const { data, isLoading, error, refetch } = useGetUsersQuery({
        page: page + 1,
        limit: pageSize,
        search,
        role: role || undefined,
        is_active: onlyActive ? true : undefined,
    });
    const [loadUserDetail, { isFetching: isLoadingDetail }] = useLazyGetUserByIdQuery();
    const [createUser, { isLoading: isCreating }] = useCreateUserMutation();
    const [updateUser, { isLoading: isUpdating }] = useUpdateUserMutation();
    const [resetUserPassword, { isLoading: isResetting }] = useResetUserPasswordMutation();

    const rows = data?.data?.users || [];

    const columns = useMemo<GridColDef[]>(() => [
        {
            field: 'full_name',
            headerName: 'Người dùng',
            minWidth: 220,
            flex: 1.2,
            renderCell: (params: GridRenderCellParams) => (
                <Stack spacing={0.25} sx={{ py: 1 }}>
                    <Typography variant="body2" sx={{ fontWeight: 700 }}>
                        {params.row.full_name}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                        {params.row.email}
                    </Typography>
                </Stack>
            ),
        },
        {
            field: 'code',
            headerName: 'Mã',
            width: 150,
        },
        {
            field: 'role',
            headerName: 'Vai trò',
            width: 140,
            renderCell: (params: GridRenderCellParams) => (
                <Chip label={params.value as string} size="small" color="primary" variant="outlined" />
            ),
        },
        {
            field: 'is_active',
            headerName: 'Trạng thái',
            width: 140,
            renderCell: (params: GridRenderCellParams) => (
                <Chip
                    label={params.value ? 'Đang hoạt động' : 'Đã khóa'}
                    size="small"
                    color={params.value ? 'success' : 'default'}
                />
            ),
        },
        {
            field: 'actions',
            headerName: 'Thao tác',
            width: 190,
            sortable: false,
            renderCell: (params: GridRenderCellParams) => (
                <Stack direction="row" spacing={1}>
                    <Button
                        size="small"
                        startIcon={<Visibility />}
                        onClick={async () => {
                            setSelectedUser(params.row);
                            await loadUserDetail(params.row.id);
                            setDetailOpen(true);
                        }}
                    >
                        Xem
                    </Button>
                    <Button
                        size="small"
                        startIcon={<Edit />}
                        onClick={() => {
                            setSelectedUser(params.row);
                            setForm({
                                code: params.row.code,
                                full_name: params.row.full_name,
                                email: params.row.email,
                                role: params.row.role,
                                is_active: params.row.is_active,
                                password: '',
                            });
                            setDialogOpen(true);
                        }}
                    >
                        Sửa
                    </Button>
                    <Button
                        size="small"
                        color="warning"
                        startIcon={<LockReset />}
                        onClick={() => {
                            setSelectedUser(params.row);
                            setResetPasswordValue('');
                            setResetDialogOpen(true);
                        }}
                    >
                        Reset
                    </Button>
                </Stack>
            ),
        },
    ], [loadUserDetail]);

    const handleSave = async () => {
        try {
            if (selectedUser) {
                await updateUser({
                    id: selectedUser.id,
                    body: {
                        full_name: form.full_name,
                        role: form.role,
                        is_active: form.is_active,
                    },
                }).unwrap();
                toast.success('Cập nhật tài khoản thành công');
            } else {
                await createUser({
                    code: form.code || undefined,
                    full_name: form.full_name,
                    email: form.email,
                    role: form.role,
                    is_active: form.is_active,
                    password: form.password,
                }).unwrap();
                toast.success('Tạo tài khoản thành công');
            }
            setDialogOpen(false);
            setSelectedUser(null);
            setForm(emptyForm);
            refetch();
        } catch (err) {
            toast.error(getApiErrorMessage(err, 'Không lưu được tài khoản.'));
        }
    };

    const handleResetPassword = async () => {
        if (!selectedUser) return;

        try {
            await resetUserPassword({
                id: selectedUser.id,
                new_password: resetPasswordValue,
            }).unwrap();
            toast.success('Đặt lại mật khẩu thành công');
            setResetDialogOpen(false);
            setResetPasswordValue('');
        } catch (err) {
            toast.error(getApiErrorMessage(err, 'Không đặt lại được mật khẩu.'));
        }
    };

    return (
        <Box>
            <PageHeader
                title="Quản lý tài khoản"
                subtitle="Admin và Super Admin quản lý user, role và trạng thái tài khoản nội bộ."
                breadcrumbs={[
                    { label: 'Dashboard', path: '/app/admin/overview' },
                    { label: 'Quản lý tài khoản' },
                ]}
                actions={(
                    <Button
                        variant="contained"
                        startIcon={<Add />}
                        onClick={() => {
                            setSelectedUser(null);
                            setForm(emptyForm);
                            setDialogOpen(true);
                        }}
                    >
                        Tạo tài khoản
                    </Button>
                )}
            />

            {error && (
                <Alert severity="error" sx={{ mb: 3 }}>
                    {getApiErrorMessage(error, 'Không tải được danh sách tài khoản.')}
                </Alert>
            )}

            <Paper elevation={0} sx={{ p: 2, borderRadius: 3, border: '1px solid #e2e8f0', mb: 3 }}>
                <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
                    <TextField
                        label="Tìm theo tên / email / mã"
                        size="small"
                        value={searchInput}
                        onChange={(event) => setSearchInput(event.target.value)}
                        onKeyDown={(event) => {
                            if (event.key === 'Enter') {
                                setPage(0);
                                setSearch(searchInput.trim());
                            }
                        }}
                        InputProps={{ startAdornment: <SearchRounded sx={{ mr: 1, color: 'text.secondary' }} /> }}
                        sx={{ minWidth: 280 }}
                    />
                    <FormControl size="small" sx={{ minWidth: 180 }}>
                        <InputLabel>Vai trò</InputLabel>
                        <Select
                            label="Vai trò"
                            value={role}
                            onChange={(event) => {
                                setPage(0);
                                setRole(event.target.value);
                            }}
                        >
                            <MenuItem value="">Tất cả</MenuItem>
                            {roles.map((item) => (
                                <MenuItem key={item} value={item}>{item}</MenuItem>
                            ))}
                        </Select>
                    </FormControl>
                    <Stack direction="row" spacing={1} alignItems="center">
                        <Switch
                            checked={onlyActive}
                            onChange={(_, checked) => {
                                setPage(0);
                                setOnlyActive(checked);
                            }}
                        />
                        <Typography variant="body2">Chỉ tài khoản đang hoạt động</Typography>
                    </Stack>
                    <Button
                        variant="outlined"
                        onClick={() => {
                            setPage(0);
                            setSearch(searchInput.trim());
                        }}
                    >
                        Áp dụng
                    </Button>
                </Stack>
            </Paper>

            <Paper elevation={0} sx={{ p: 2, borderRadius: 3, border: '1px solid #e2e8f0' }}>
                <DataGrid
                    rows={rows}
                    columns={columns}
                    loading={isLoading}
                    paginationMode="server"
                    rowCount={data?.data?.pagination?.total_items || 0}
                    paginationModel={{ page, pageSize }}
                    onPaginationModelChange={(model) => {
                        setPage(model.page);
                        setPageSize(model.pageSize);
                    }}
                    pageSizeOptions={[10, 25, 50]}
                    disableRowSelectionOnClick
                    autoHeight
                    getRowId={(row) => row.id}
                    sx={{
                        border: 'none',
                        '& .MuiDataGrid-columnHeaders': {
                            backgroundColor: '#f8fafc',
                        },
                    }}
                />
            </Paper>

            <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)} maxWidth="sm" fullWidth>
                <DialogTitle>{selectedUser ? 'Cập nhật tài khoản' : 'Tạo tài khoản mới'}</DialogTitle>
                <DialogContent>
                    <Stack spacing={2} sx={{ mt: 1 }}>
                        <TextField
                            label="Mã tài khoản"
                            value={form.code}
                            onChange={(event) => setForm((current) => ({ ...current, code: event.target.value }))}
                            disabled={!!selectedUser}
                        />
                        <TextField
                            label="Họ và tên"
                            value={form.full_name}
                            onChange={(event) => setForm((current) => ({ ...current, full_name: event.target.value }))}
                        />
                        <TextField
                            label="Email"
                            type="email"
                            value={form.email}
                            onChange={(event) => setForm((current) => ({ ...current, email: event.target.value }))}
                            disabled={!!selectedUser}
                        />
                        <FormControl fullWidth>
                            <InputLabel>Vai trò</InputLabel>
                            <Select
                                label="Vai trò"
                                value={form.role}
                                onChange={(event) => setForm((current) => ({ ...current, role: event.target.value }))}
                            >
                                {roles.map((item) => (
                                    <MenuItem key={item} value={item}>{item}</MenuItem>
                                ))}
                            </Select>
                        </FormControl>
                        {!selectedUser && (
                            <TextField
                                label="Mật khẩu khởi tạo"
                                type="password"
                                value={form.password}
                                onChange={(event) => setForm((current) => ({ ...current, password: event.target.value }))}
                            />
                        )}
                        <Stack direction="row" spacing={1} alignItems="center">
                            <Switch
                                checked={form.is_active}
                                onChange={(_, checked) => setForm((current) => ({ ...current, is_active: checked }))}
                            />
                            <Typography variant="body2">Tài khoản đang hoạt động</Typography>
                        </Stack>
                    </Stack>
                </DialogContent>
                <DialogActions sx={{ p: 2 }}>
                    <Button onClick={() => setDialogOpen(false)}>Hủy</Button>
                    <Button
                        variant="contained"
                        onClick={handleSave}
                        disabled={isCreating || isUpdating}
                    >
                        {isCreating || isUpdating ? 'Đang lưu...' : 'Lưu'}
                    </Button>
                </DialogActions>
            </Dialog>

            <Dialog open={detailOpen} onClose={() => setDetailOpen(false)} maxWidth="sm" fullWidth>
                <DialogTitle>Chi tiết tài khoản</DialogTitle>
                <DialogContent>
                    {selectedUser && (
                        <Stack spacing={1.5} sx={{ mt: 1 }}>
                            <Typography><strong>Họ tên:</strong> {selectedUser.full_name}</Typography>
                            <Typography><strong>Email:</strong> {selectedUser.email}</Typography>
                            <Typography><strong>Mã:</strong> {selectedUser.code}</Typography>
                            <Typography><strong>Vai trò:</strong> {selectedUser.role}</Typography>
                            <Typography><strong>Trạng thái:</strong> {selectedUser.is_active ? 'Đang hoạt động' : 'Đã khóa'}</Typography>
                            <Typography><strong>Tạo lúc:</strong> {selectedUser.created_at ? new Date(selectedUser.created_at).toLocaleString('vi-VN') : 'N/A'}</Typography>
                            <Typography><strong>Cập nhật lúc:</strong> {selectedUser.updated_at ? new Date(selectedUser.updated_at).toLocaleString('vi-VN') : 'N/A'}</Typography>
                            {isLoadingDetail && <Typography color="text.secondary">Đang tải thêm dữ liệu...</Typography>}
                        </Stack>
                    )}
                </DialogContent>
                <DialogActions sx={{ p: 2 }}>
                    <Button onClick={() => setDetailOpen(false)}>Đóng</Button>
                </DialogActions>
            </Dialog>

            <ConfirmDialog
                open={resetDialogOpen}
                title="Đặt lại mật khẩu"
                message={selectedUser ? `Bạn đang đặt lại mật khẩu cho tài khoản ${selectedUser.email}. Mật khẩu mới sẽ được áp dụng ngay.` : ''}
                onClose={() => setResetDialogOpen(false)}
                onConfirm={handleResetPassword}
                confirmText="Xác nhận reset"
                loading={isResetting}
            />

            {resetDialogOpen && (
                <Dialog open={resetDialogOpen} onClose={() => setResetDialogOpen(false)} maxWidth="xs" fullWidth>
                    <DialogTitle>Mật khẩu mới</DialogTitle>
                    <DialogContent>
                        <TextField
                            autoFocus
                            margin="dense"
                            label="Nhập mật khẩu mới"
                            type="password"
                            fullWidth
                            value={resetPasswordValue}
                            onChange={(event) => setResetPasswordValue(event.target.value)}
                        />
                    </DialogContent>
                    <DialogActions sx={{ p: 2 }}>
                        <Button onClick={() => setResetDialogOpen(false)}>Hủy</Button>
                        <Button variant="contained" onClick={handleResetPassword} disabled={isResetting}>
                            {isResetting ? 'Đang reset...' : 'Reset mật khẩu'}
                        </Button>
                    </DialogActions>
                </Dialog>
            )}
        </Box>
    );
};
