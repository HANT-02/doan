import React, { useMemo } from 'react';
import { AlertTriangle, BookOpen, Calendar, GraduationCap, Star } from 'lucide-react';
import { addDays, format, isSameDay, parseISO } from 'date-fns';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
    type StudentAtRiskPrediction,
    useGetStudentAcademicRecordsQuery,
    useGetStudentAttendanceQuery,
    useGetStudentAtRiskPredictionQuery,
    useGetStudentTimetableQuery,
} from '@/api/studentPortalApi';

function buildRecommendations(prediction: StudentAtRiskPrediction | null) {
    if (!prediction) {
        return ['Tiếp tục duy trì chuyên cần và hoàn thành bài tập đúng hạn.'];
    }

    const actions: string[] = [];
    if (prediction.feature_summary.attendance_rate_28d < 0.8) {
        actions.push('Cải thiện chuyên cần trong 2-4 tuần tới và hạn chế nghỉ học không cần thiết.');
    }
    if (prediction.feature_summary.homework_completion_rate_28d < 0.8) {
        actions.push('Ưu tiên hoàn thành bài tập và nộp đúng hạn ở các buổi tiếp theo.');
    }
    if (prediction.feature_summary.average_total_score_28d < 5) {
        actions.push('Trao đổi với giáo viên để được hỗ trợ thêm về phần kiến thức còn yếu.');
    }
    if (!actions.length) {
        actions.push('Giữ ổn định tiến độ học tập và tiếp tục theo dõi cảnh báo định kỳ.');
    }
    return actions.slice(0, 3);
}

export const StudentOverview: React.FC = () => {
    const { data: timetableData } = useGetStudentTimetableQuery({
        from: format(new Date(), 'yyyy-MM-dd'),
        to: format(addDays(new Date(), 7), 'yyyy-MM-dd'),
    });
    const { data: attendanceData } = useGetStudentAttendanceQuery();
    const { data: recordsData } = useGetStudentAcademicRecordsQuery();
    const { data: atRiskData } = useGetStudentAtRiskPredictionQuery();

    const lessons = timetableData?.data?.lessons ?? [];
    const attendanceSummary = attendanceData?.data?.summary;
    const classSummaries = recordsData?.data?.class_summaries ?? [];
    const prediction = atRiskData?.data?.prediction ?? null;

    const todayLessons = useMemo(
        () => lessons.filter((lesson) => isSameDay(parseISO(lesson.date_start), new Date())),
        [lessons],
    );

    const averageScore = useMemo(() => {
        if (!classSummaries.length) {
            return 0;
        }
        return classSummaries.reduce((sum, item) => sum + item.average_total_score, 0) / classSummaries.length;
    }, [classSummaries]);

    const stats = [
        { label: 'Điểm TB', value: averageScore ? averageScore.toFixed(2) : '--', icon: Star, color: 'text-yellow-500', hint: 'Tính trên các lớp đã có dữ liệu' },
        { label: 'Lớp đang học', value: String(classSummaries.length), icon: GraduationCap, color: 'text-blue-600', hint: 'Có dữ liệu kết quả học tập' },
        { label: 'Buổi học hôm nay', value: String(todayLessons.length), icon: BookOpen, color: 'text-green-600', hint: 'Dựa trên lịch học hiện tại' },
        {
            label: 'Chuyên cần',
            value: attendanceSummary ? `${Math.round(attendanceSummary.attendance_rate * 100)}%` : '--',
            icon: Calendar,
            color: attendanceSummary?.warning ? 'text-orange-600' : 'text-purple-600',
            hint: 'Tổng hợp 28 ngày gần nhất',
        },
    ];

    const recommendations = buildRecommendations(prediction);

    return (
        <div className="space-y-6 animate-in fade-in-50 duration-500">
            <div className="space-y-1">
                <h1 className="text-3xl font-bold tracking-tight text-gray-900">Cổng học sinh</h1>
                <p className="text-sm text-gray-500">Tổng quan nhanh về lịch học, chuyên cần, kết quả học tập và cảnh báo sớm.</p>
            </div>

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                {stats.map((stat) => {
                    const Icon = stat.icon;
                    return (
                        <Card key={stat.label}>
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">{stat.label}</CardTitle>
                                <Icon className={`h-4 w-4 ${stat.color}`} />
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold">{stat.value}</div>
                                <p className="text-xs text-muted-foreground">{stat.hint}</p>
                            </CardContent>
                        </Card>
                    );
                })}
            </div>

            {prediction ? (
                <Card className={prediction.risk_label === 'AT_RISK' ? 'border-orange-300 bg-orange-50/40' : ''}>
                    <CardHeader className="flex flex-row items-start justify-between space-y-0">
                        <div className="space-y-1">
                            <CardTitle className="text-lg">Cảnh báo học tập AT_RISK</CardTitle>
                            <p className="text-sm text-muted-foreground">
                                Mô hình {prediction.model_name} • mức {prediction.risk_band.toLowerCase()} • xác suất {Math.round(prediction.risk_score * 100)}%
                            </p>
                        </div>
                        <AlertTriangle className={`h-5 w-5 ${prediction.risk_label === 'AT_RISK' ? 'text-orange-600' : 'text-emerald-600'}`} />
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="rounded-md border bg-background p-3">
                            <p className="text-sm font-medium">{prediction.primary_reason}</p>
                            <p className="mt-1 text-xs text-muted-foreground">
                                Lớp {prediction.class_name} ({prediction.class_code}) • snapshot {new Date(prediction.snapshot_at).toLocaleString('vi-VN')}
                            </p>
                        </div>

                        <div className="grid gap-3 md:grid-cols-3">
                            {prediction.top_features.map((item) => (
                                <div key={item.key} className="rounded-md border bg-background p-3">
                                    <p className="text-xs uppercase tracking-wide text-muted-foreground">{item.label}</p>
                                    <p className="mt-1 text-lg font-semibold">{item.display_value}</p>
                                </div>
                            ))}
                        </div>

                        <div className="space-y-2">
                            <p className="text-sm font-semibold">Gợi ý hành động</p>
                            <ul className="list-disc space-y-1 pl-5 text-sm text-muted-foreground">
                                {recommendations.map((item) => (
                                    <li key={item}>{item}</li>
                                ))}
                            </ul>
                        </div>
                    </CardContent>
                </Card>
            ) : (
                <Card>
                    <CardHeader>
                        <CardTitle className="text-lg">Cảnh báo học tập</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p className="text-sm text-muted-foreground">
                            Chưa có dữ liệu dự báo AT_RISK tại thời điểm này. Hệ thống sẽ cập nhật khi đủ dữ liệu học tập và chuyên cần.
                        </p>
                    </CardContent>
                </Card>
            )}

            <Card>
                <CardHeader>
                    <CardTitle className="text-lg">Lịch học hôm nay</CardTitle>
                </CardHeader>
                <CardContent>
                    {todayLessons.length ? (
                        <div className="space-y-3">
                            {todayLessons.map((lesson) => (
                                <div key={lesson.id} className="rounded-md border p-3">
                                    <div className="flex flex-wrap items-center justify-between gap-2">
                                        <div>
                                            <p className="font-medium">{lesson.class_name}</p>
                                            <p className="text-sm text-muted-foreground">
                                                {lesson.class_code} • {lesson.shift?.name || 'Chưa gắn ca học'} • {lesson.room_name || 'Chưa xếp phòng'}
                                            </p>
                                        </div>
                                        <span className="text-sm text-muted-foreground">
                                            {new Date(lesson.date_start).toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' })} - {new Date(lesson.date_end).toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' })}
                                        </span>
                                    </div>
                                </div>
                            ))}
                        </div>
                    ) : (
                        <div className="text-sm text-gray-500">Hôm nay bạn chưa có buổi học nào theo lịch hiện tại.</div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
};
