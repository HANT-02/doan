import { useState } from 'react';
import {
    Box,
    TextField,
    Select,
    MenuItem,
    FormControl,
    InputLabel,
    Button,
    Grid,
    Paper,
} from '@mui/material';
import { Search, FilterAlt, Clear } from '@mui/icons-material';

export interface TeacherFiltersState {
    search: string;
    status: string;
    employment_type: string;
}

interface TeacherFiltersProps {
    filters: TeacherFiltersState;
    onChange: (filters: TeacherFiltersState) => void;
}

export const TeacherFilters = ({ filters, onChange }: TeacherFiltersProps) => {
    const [localFilters, setLocalFilters] = useState<TeacherFiltersState>(filters);

    const handleChange = (field: keyof TeacherFiltersState, value: string) => {
        const newFilters = { ...localFilters, [field]: value };
        setLocalFilters(newFilters);
        onChange(newFilters);
    };

    const handleReset = () => {
        const resetFilters: TeacherFiltersState = {
            search: '',
            status: '',
            employment_type: '',
        };
        setLocalFilters(resetFilters);
        onChange(resetFilters);
    };

    const hasActiveFilters =
        localFilters.search || localFilters.status || localFilters.employment_type;

    return (
        <Paper sx={{ p: 2, mb: 3 }}>
            <Grid container spacing={2} alignItems="center">
                <Grid item xs={12} md={4}>
                    <TextField
                        fullWidth
                        size="small"
                        label="Tìm kiếm"
                        placeholder="Tên, email, điện thoại, mã..."
                        value={localFilters.search}
                        onChange={(e) => handleChange('search', e.target.value)}
                        InputProps={{
                            startAdornment: <Search sx={{ mr: 1, color: 'text.secondary' }} />,
                        }}
                    />
                </Grid>

                <Grid item xs={12} sm={6} md={3}>
                    <FormControl fullWidth size="small">
                        <InputLabel>Trạng thái</InputLabel>
                        <Select
                            value={localFilters.status}
                            label="Trạng thái"
                            onChange={(e) => handleChange('status', e.target.value)}
                        >
                            <MenuItem value="">Tất cả</MenuItem>
                            <MenuItem value="ACTIVE">Hoạt động</MenuItem>
                            <MenuItem value="INACTIVE">Không hoạt động</MenuItem>
                        </Select>
                    </FormControl>
                </Grid>

                <Grid item xs={12} sm={6} md={3}>
                    <FormControl fullWidth size="small">
                        <InputLabel>Loại hình làm việc</InputLabel>
                        <Select
                            value={localFilters.employment_type}
                            label="Loại hình làm việc"
                            onChange={(e) => handleChange('employment_type', e.target.value)}
                        >
                            <MenuItem value="">Tất cả</MenuItem>
                            <MenuItem value="FULL_TIME">Toàn thời gian</MenuItem>
                            <MenuItem value="PART_TIME">Bán thời gian</MenuItem>
                        </Select>
                    </FormControl>
                </Grid>

                <Grid item xs={12} md={2}>
                    <Button
                        fullWidth
                        variant="outlined"
                        startIcon={hasActiveFilters ? <Clear /> : <FilterAlt />}
                        onClick={handleReset}
                        disabled={!hasActiveFilters}
                    >
                        {hasActiveFilters ? 'Xóa bộ lọc' : 'Bộ lọc'}
                    </Button>
                </Grid>
            </Grid>
        </Paper>
    );
};
