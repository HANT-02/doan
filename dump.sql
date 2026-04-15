--
-- PostgreSQL database dump
--

\restrict DdX62PsYk80fV80n3jObLeAKtAWXehdJI5e0N4kS5jtdnMd5yqMXkgP3cQOg3w4

-- Dumped from database version 16.11 (Debian 16.11-1.pgdg13+1)
-- Dumped by pg_dump version 16.11 (Debian 16.11-1.pgdg13+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE IF EXISTS ONLY public.programs DROP CONSTRAINT IF EXISTS fk_programs_created_by;
ALTER TABLE IF EXISTS ONLY public.programs DROP CONSTRAINT IF EXISTS fk_programs_approved_by;
ALTER TABLE IF EXISTS ONLY public.program_courses DROP CONSTRAINT IF EXISTS fk_program_courses_program;
ALTER TABLE IF EXISTS ONLY public.program_courses DROP CONSTRAINT IF EXISTS fk_program_courses_course;
ALTER TABLE IF EXISTS ONLY public.outcomes DROP CONSTRAINT IF EXISTS fk_outcomes_program;
ALTER TABLE IF EXISTS ONLY public.outcomes DROP CONSTRAINT IF EXISTS fk_outcomes_objective;
ALTER TABLE IF EXISTS ONLY public.objectives DROP CONSTRAINT IF EXISTS fk_objectives_program;
ALTER TABLE IF EXISTS ONLY public.lessons DROP CONSTRAINT IF EXISTS fk_lessons_teacher;
ALTER TABLE IF EXISTS ONLY public.lessons DROP CONSTRAINT IF EXISTS fk_lessons_room;
ALTER TABLE IF EXISTS ONLY public.lessons DROP CONSTRAINT IF EXISTS fk_lessons_class;
ALTER TABLE IF EXISTS ONLY public.lesson_summaries DROP CONSTRAINT IF EXISTS fk_lesson_summaries_lesson;
ALTER TABLE IF EXISTS ONLY public.lesson_summaries DROP CONSTRAINT IF EXISTS fk_lesson_summaries_created_by;
ALTER TABLE IF EXISTS ONLY public.leave_requests DROP CONSTRAINT IF EXISTS fk_leave_requests_student;
ALTER TABLE IF EXISTS ONLY public.leave_requests DROP CONSTRAINT IF EXISTS fk_leave_requests_lesson;
ALTER TABLE IF EXISTS ONLY public.leave_requests DROP CONSTRAINT IF EXISTS fk_leave_requests_class;
ALTER TABLE IF EXISTS ONLY public.leave_requests DROP CONSTRAINT IF EXISTS fk_leave_requests_approved_by;
ALTER TABLE IF EXISTS ONLY public.enrollments DROP CONSTRAINT IF EXISTS fk_enrollments_student;
ALTER TABLE IF EXISTS ONLY public.enrollments DROP CONSTRAINT IF EXISTS fk_enrollments_class;
ALTER TABLE IF EXISTS ONLY public.classes DROP CONSTRAINT IF EXISTS fk_classes_teacher;
ALTER TABLE IF EXISTS ONLY public.classes DROP CONSTRAINT IF EXISTS fk_classes_room;
ALTER TABLE IF EXISTS ONLY public.classes DROP CONSTRAINT IF EXISTS fk_classes_program;
ALTER TABLE IF EXISTS ONLY public.classes DROP CONSTRAINT IF EXISTS fk_classes_course;
ALTER TABLE IF EXISTS ONLY public.class_schedules DROP CONSTRAINT IF EXISTS fk_classes_class_schedules;
ALTER TABLE IF EXISTS ONLY public.class_schedules DROP CONSTRAINT IF EXISTS fk_class_schedules_shift;
ALTER TABLE IF EXISTS ONLY public.class_schedules DROP CONSTRAINT IF EXISTS fk_class_schedules_room;
ALTER TABLE IF EXISTS ONLY public.class_schedules DROP CONSTRAINT IF EXISTS fk_class_schedules_class;
ALTER TABLE IF EXISTS ONLY public.attendances DROP CONSTRAINT IF EXISTS fk_attendances_student;
ALTER TABLE IF EXISTS ONLY public.attendances DROP CONSTRAINT IF EXISTS fk_attendances_lesson;
ALTER TABLE IF EXISTS ONLY public.academic_records DROP CONSTRAINT IF EXISTS fk_academic_records_student;
ALTER TABLE IF EXISTS ONLY public.academic_records DROP CONSTRAINT IF EXISTS fk_academic_records_lesson_summary;
DROP INDEX IF EXISTS public.idx_users_deleted_at;
DROP INDEX IF EXISTS public.idx_user_otps_user_id;
DROP INDEX IF EXISTS public.idx_user_otps_deleted_at;
DROP INDEX IF EXISTS public.idx_teachers_deleted_at;
DROP INDEX IF EXISTS public.idx_students_deleted_at;
DROP INDEX IF EXISTS public.idx_shifts_deleted_at;
DROP INDEX IF EXISTS public.idx_shifts_code;
DROP INDEX IF EXISTS public.idx_programs_deleted_at;
DROP INDEX IF EXISTS public.idx_password_resets_deleted_at;
DROP INDEX IF EXISTS public.idx_courses_deleted_at;
DROP INDEX IF EXISTS public.idx_classes_deleted_at;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY public.user_otps DROP CONSTRAINT IF EXISTS user_otps_pkey;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS uni_users_email;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS uni_users_code;
ALTER TABLE IF EXISTS ONLY public.teachers DROP CONSTRAINT IF EXISTS uni_teachers_code;
ALTER TABLE IF EXISTS ONLY public.students DROP CONSTRAINT IF EXISTS uni_students_code;
ALTER TABLE IF EXISTS ONLY public.rooms DROP CONSTRAINT IF EXISTS uni_rooms_code;
ALTER TABLE IF EXISTS ONLY public.programs DROP CONSTRAINT IF EXISTS uni_programs_code;
ALTER TABLE IF EXISTS ONLY public.outcomes DROP CONSTRAINT IF EXISTS uni_outcomes_code;
ALTER TABLE IF EXISTS ONLY public.objectives DROP CONSTRAINT IF EXISTS uni_objectives_code;
ALTER TABLE IF EXISTS ONLY public.lesson_summaries DROP CONSTRAINT IF EXISTS uni_lesson_summaries_lesson_id;
ALTER TABLE IF EXISTS ONLY public.courses DROP CONSTRAINT IF EXISTS uni_courses_code;
ALTER TABLE IF EXISTS ONLY public.classes DROP CONSTRAINT IF EXISTS uni_classes_code;
ALTER TABLE IF EXISTS ONLY public.teachers DROP CONSTRAINT IF EXISTS teachers_pkey;
ALTER TABLE IF EXISTS ONLY public.students DROP CONSTRAINT IF EXISTS students_pkey;
ALTER TABLE IF EXISTS ONLY public.shifts DROP CONSTRAINT IF EXISTS shifts_pkey;
ALTER TABLE IF EXISTS ONLY public.schema_migrations DROP CONSTRAINT IF EXISTS schema_migrations_pkey;
ALTER TABLE IF EXISTS ONLY public.rooms DROP CONSTRAINT IF EXISTS rooms_pkey;
ALTER TABLE IF EXISTS ONLY public.programs DROP CONSTRAINT IF EXISTS programs_pkey;
ALTER TABLE IF EXISTS ONLY public.program_courses DROP CONSTRAINT IF EXISTS program_courses_pkey;
ALTER TABLE IF EXISTS ONLY public.password_resets DROP CONSTRAINT IF EXISTS password_resets_pkey;
ALTER TABLE IF EXISTS ONLY public.outcomes DROP CONSTRAINT IF EXISTS outcomes_pkey;
ALTER TABLE IF EXISTS ONLY public.objectives DROP CONSTRAINT IF EXISTS objectives_pkey;
ALTER TABLE IF EXISTS ONLY public.lessons DROP CONSTRAINT IF EXISTS lessons_pkey;
ALTER TABLE IF EXISTS ONLY public.lesson_summaries DROP CONSTRAINT IF EXISTS lesson_summaries_pkey;
ALTER TABLE IF EXISTS ONLY public.leave_requests DROP CONSTRAINT IF EXISTS leave_requests_pkey;
ALTER TABLE IF EXISTS ONLY public.enrollments DROP CONSTRAINT IF EXISTS enrollments_pkey;
ALTER TABLE IF EXISTS ONLY public.courses DROP CONSTRAINT IF EXISTS courses_pkey;
ALTER TABLE IF EXISTS ONLY public.consultations DROP CONSTRAINT IF EXISTS consultations_pkey;
ALTER TABLE IF EXISTS ONLY public.classes DROP CONSTRAINT IF EXISTS classes_pkey;
ALTER TABLE IF EXISTS ONLY public.class_schedules DROP CONSTRAINT IF EXISTS class_schedules_pkey;
ALTER TABLE IF EXISTS ONLY public.attendances DROP CONSTRAINT IF EXISTS attendances_pkey;
ALTER TABLE IF EXISTS ONLY public.academic_records DROP CONSTRAINT IF EXISTS academic_records_pkey;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.user_otps;
DROP TABLE IF EXISTS public.teachers;
DROP TABLE IF EXISTS public.students;
DROP TABLE IF EXISTS public.shifts;
DROP TABLE IF EXISTS public.schema_migrations;
DROP TABLE IF EXISTS public.rooms;
DROP TABLE IF EXISTS public.programs;
DROP TABLE IF EXISTS public.program_courses;
DROP TABLE IF EXISTS public.password_resets;
DROP TABLE IF EXISTS public.outcomes;
DROP TABLE IF EXISTS public.objectives;
DROP TABLE IF EXISTS public.lessons;
DROP TABLE IF EXISTS public.lesson_summaries;
DROP TABLE IF EXISTS public.leave_requests;
DROP TABLE IF EXISTS public.enrollments;
DROP TABLE IF EXISTS public.courses;
DROP TABLE IF EXISTS public.consultations;
DROP TABLE IF EXISTS public.classes;
DROP TABLE IF EXISTS public.class_schedules;
DROP TABLE IF EXISTS public.attendances;
DROP TABLE IF EXISTS public.academic_records;
DROP EXTENSION IF EXISTS "uuid-ossp";
--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: academic_records; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.academic_records (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    lesson_summary_id uuid NOT NULL,
    student_id uuid NOT NULL,
    homework_completed boolean DEFAULT false,
    homework_score numeric(5,2),
    attitude_rating bigint,
    participation_score numeric(5,2),
    personal_comment text,
    total_score numeric(5,2),
    is_completed boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone
);


ALTER TABLE public.academic_records OWNER TO root;

--
-- Name: attendances; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.attendances (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    lesson_id uuid NOT NULL,
    student_id uuid NOT NULL,
    status bigint NOT NULL,
    note text,
    marked_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone
);


ALTER TABLE public.attendances OWNER TO root;

--
-- Name: class_schedules; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.class_schedules (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    class_id uuid NOT NULL,
    day_of_week character varying(20) NOT NULL,
    start_time character varying(10) NOT NULL,
    end_time character varying(10) NOT NULL,
    room_id uuid,
    shift_id uuid NOT NULL
);


ALTER TABLE public.class_schedules OWNER TO root;

--
-- Name: classes; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.classes (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    code character varying(50) NOT NULL,
    name character varying(255) NOT NULL,
    notes text,
    start_date timestamp with time zone NOT NULL,
    end_date timestamp with time zone,
    max_students bigint,
    status character varying(50) DEFAULT 'OPEN'::character varying,
    price numeric(10,2),
    program_id uuid,
    course_id uuid,
    teacher_id uuid,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    room_id uuid
);


ALTER TABLE public.classes OWNER TO root;

--
-- Name: consultations; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.consultations (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    full_name character varying(255) NOT NULL,
    phone character varying(20) NOT NULL,
    grade_level character varying(50) NOT NULL,
    notes text,
    status character varying(50) DEFAULT 'PENDING'::character varying,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone
);


ALTER TABLE public.consultations OWNER TO root;

--
-- Name: courses; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.courses (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    code character varying(50) NOT NULL,
    name character varying(255),
    description text,
    grade_level character varying(50),
    subject character varying(255),
    session_count bigint,
    session_duration_minutes bigint,
    total_hours numeric(8,2),
    price numeric(10,2),
    status character varying(50) DEFAULT 'ACTIVE'::character varying,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.courses OWNER TO root;

--
-- Name: enrollments; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.enrollments (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    class_id uuid NOT NULL,
    student_id uuid NOT NULL,
    status character varying(50) DEFAULT 'APPLIED'::character varying,
    approved_at timestamp with time zone,
    rejected_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone
);


ALTER TABLE public.enrollments OWNER TO root;

--
-- Name: leave_requests; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.leave_requests (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    student_id uuid NOT NULL,
    leave_type character varying(50) NOT NULL,
    apply_date timestamp with time zone NOT NULL,
    late_minutes bigint,
    early_minutes bigint,
    reason text NOT NULL,
    documents text[],
    class_id uuid,
    lesson_id uuid,
    subject character varying(255),
    status character varying(50) DEFAULT 'PENDING'::character varying,
    approved_by_id uuid,
    approved_at timestamp with time zone,
    rejection_reason text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone
);


ALTER TABLE public.leave_requests OWNER TO root;

--
-- Name: lesson_summaries; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.lesson_summaries (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    lesson_id uuid NOT NULL,
    topic text,
    lesson_content text,
    class_feedback text,
    homework text,
    homework_deadline timestamp with time zone,
    teacher_notes text,
    created_by_id uuid,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone
);


ALTER TABLE public.lesson_summaries OWNER TO root;

--
-- Name: lessons; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.lessons (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    class_id uuid NOT NULL,
    date_start timestamp with time zone NOT NULL,
    date_end timestamp with time zone NOT NULL,
    room_id uuid,
    notes text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    teacher_id uuid
);


ALTER TABLE public.lessons OWNER TO root;

--
-- Name: objectives; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.objectives (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    code character varying(50) NOT NULL,
    name text NOT NULL,
    program_id uuid NOT NULL
);


ALTER TABLE public.objectives OWNER TO root;

--
-- Name: outcomes; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.outcomes (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    code character varying(50) NOT NULL,
    name text NOT NULL,
    program_id uuid NOT NULL,
    objective_id uuid
);


ALTER TABLE public.outcomes OWNER TO root;

--
-- Name: password_resets; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.password_resets (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    token_hash text NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    used_at timestamp with time zone,
    requested_ip character varying(45),
    user_agent text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone
);


ALTER TABLE public.password_resets OWNER TO root;

--
-- Name: program_courses; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.program_courses (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    program_id uuid NOT NULL,
    course_id uuid NOT NULL
);


ALTER TABLE public.program_courses OWNER TO root;

--
-- Name: programs; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.programs (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    code character varying(50) NOT NULL,
    name character varying(255),
    track character varying(50),
    effective_from timestamp with time zone,
    effective_to timestamp with time zone,
    created_by_id uuid,
    approved_by_id uuid,
    approval_note text,
    published_at timestamp with time zone,
    archived_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    status character varying(50) DEFAULT 'ACTIVE'::character varying
);


ALTER TABLE public.programs OWNER TO root;

--
-- Name: rooms; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.rooms (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    code character varying(50) NOT NULL,
    name character varying(255) NOT NULL,
    capacity bigint,
    address text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.rooms OWNER TO root;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO root;

--
-- Name: shifts; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.shifts (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    code character varying(50) NOT NULL,
    name character varying(255) NOT NULL,
    start_time character varying(10) NOT NULL,
    end_time character varying(10) NOT NULL,
    duration_minutes bigint NOT NULL,
    session_type character varying(50) NOT NULL,
    is_active boolean DEFAULT true,
    notes text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.shifts OWNER TO root;

--
-- Name: students; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.students (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    code character varying(50),
    full_name character varying(255),
    email character varying(255),
    phone character varying(20),
    guardian_phone character varying(20),
    grade_level character varying(50),
    status character varying(50) DEFAULT 'ACTIVE'::character varying,
    date_of_birth timestamp with time zone,
    gender character varying(20),
    address text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.students OWNER TO root;

--
-- Name: teachers; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.teachers (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    code character varying(50),
    full_name character varying(255),
    email character varying(255),
    phone character varying(20),
    is_school_teacher boolean DEFAULT false,
    school_name character varying(255),
    employment_type character varying(50) DEFAULT 'PART_TIME'::character varying,
    status character varying(50) DEFAULT 'ACTIVE'::character varying,
    notes text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.teachers OWNER TO root;

--
-- Name: user_otps; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.user_otps (
    id text NOT NULL,
    user_id uuid NOT NULL,
    otp_hash text,
    expired_at timestamp with time zone,
    used_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone
);


ALTER TABLE public.user_otps OWNER TO root;

--
-- Name: users; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    code character varying(50),
    full_name character varying(255),
    email character varying(255) NOT NULL,
    password text,
    role character varying(50) DEFAULT 'STUDENT'::character varying,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.users OWNER TO root;

--
-- Data for Name: academic_records; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.academic_records (id, lesson_summary_id, student_id, homework_completed, homework_score, attitude_rating, participation_score, personal_comment, total_score, is_completed, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: attendances; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.attendances (id, lesson_id, student_id, status, note, marked_at, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: class_schedules; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.class_schedules (id, class_id, day_of_week, start_time, end_time, room_id, shift_id) FROM stdin;
\.


--
-- Data for Name: classes; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.classes (id, code, name, notes, start_date, end_date, max_students, status, price, program_id, course_id, teacher_id, created_at, updated_at, deleted_at, room_id) FROM stdin;
66666666-6666-6666-6666-000000000001	L-001	Lớp Học 1	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000001	22222222-2222-2222-2222-000000000001	44444444-4444-4444-4444-000000000001	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000003	L-003	Lớp Học 3	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000003	22222222-2222-2222-2222-000000000003	44444444-4444-4444-4444-000000000003	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000004	L-004	Lớp Học 4	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000004	22222222-2222-2222-2222-000000000004	44444444-4444-4444-4444-000000000004	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000005	L-005	Lớp Học 5	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000005	22222222-2222-2222-2222-000000000005	44444444-4444-4444-4444-000000000005	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000006	L-006	Lớp Học 6	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000006	22222222-2222-2222-2222-000000000006	44444444-4444-4444-4444-000000000006	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000007	L-007	Lớp Học 7	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000007	22222222-2222-2222-2222-000000000007	44444444-4444-4444-4444-000000000007	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000008	L-008	Lớp Học 8	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000008	22222222-2222-2222-2222-000000000008	44444444-4444-4444-4444-000000000008	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000009	L-009	Lớp Học 9	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000009	22222222-2222-2222-2222-000000000009	44444444-4444-4444-4444-000000000009	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000010	L-010	Lớp Học 10	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000010	22222222-2222-2222-2222-000000000010	44444444-4444-4444-4444-000000000010	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000011	L-011	Lớp Học 11	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000011	22222222-2222-2222-2222-000000000011	44444444-4444-4444-4444-000000000011	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000012	L-012	Lớp Học 12	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000012	22222222-2222-2222-2222-000000000012	44444444-4444-4444-4444-000000000012	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000013	L-013	Lớp Học 13	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000013	22222222-2222-2222-2222-000000000013	44444444-4444-4444-4444-000000000013	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000014	L-014	Lớp Học 14	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000014	22222222-2222-2222-2222-000000000014	44444444-4444-4444-4444-000000000014	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000015	L-015	Lớp Học 15	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000015	22222222-2222-2222-2222-000000000015	44444444-4444-4444-4444-000000000015	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000016	L-016	Lớp Học 16	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000016	22222222-2222-2222-2222-000000000016	44444444-4444-4444-4444-000000000016	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000017	L-017	Lớp Học 17	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000017	22222222-2222-2222-2222-000000000017	44444444-4444-4444-4444-000000000017	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000018	L-018	Lớp Học 18	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000018	22222222-2222-2222-2222-000000000018	44444444-4444-4444-4444-000000000018	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000019	L-019	Lớp Học 19	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000019	22222222-2222-2222-2222-000000000019	44444444-4444-4444-4444-000000000019	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000020	L-020	Lớp Học 20	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000020	22222222-2222-2222-2222-000000000020	44444444-4444-4444-4444-000000000020	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000021	L-021	Lớp Học 21	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000021	22222222-2222-2222-2222-000000000021	44444444-4444-4444-4444-000000000021	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000022	L-022	Lớp Học 22	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000022	22222222-2222-2222-2222-000000000022	44444444-4444-4444-4444-000000000022	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000023	L-023	Lớp Học 23	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000023	22222222-2222-2222-2222-000000000023	44444444-4444-4444-4444-000000000023	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000024	L-024	Lớp Học 24	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000024	22222222-2222-2222-2222-000000000024	44444444-4444-4444-4444-000000000024	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000025	L-025	Lớp Học 25	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000025	22222222-2222-2222-2222-000000000025	44444444-4444-4444-4444-000000000025	2026-02-26 03:06:23.665713+00	\N	\N	\N
66666666-6666-6666-6666-000000000002	L-002	Lớp Học 2	\N	2026-02-26 03:06:23.665713+00	\N	20	OPEN	2000000.00	33333333-3333-3333-3333-000000000002	22222222-2222-2222-2222-000000000002	44444444-4444-4444-4444-000000000003	2026-02-26 03:06:23.665713+00	2026-03-31 02:10:46.329201+00	\N	\N
\.


--
-- Data for Name: consultations; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.consultations (id, full_name, phone, grade_level, notes, status, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: courses; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.courses (id, code, name, description, grade_level, subject, session_count, session_duration_minutes, total_hours, price, status, created_at, updated_at, deleted_at) FROM stdin;
22222222-2222-2222-2222-000000000001	CRS-001	Khóa học Tiếng Anh Tăng Cường 1	Nội dung chuẩn kiến thức Tiếng Anh	Khối 9	Tiếng Anh	24	90	36.00	1510000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000002	CRS-002	Khóa học Vật Lý Tăng Cường 2	Nội dung chuẩn kiến thức Vật Lý	Khối 10	Vật Lý	24	90	36.00	1520000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000003	CRS-003	Khóa học Hóa Học Tăng Cường 3	Nội dung chuẩn kiến thức Hóa Học	Khối 11	Hóa Học	24	90	36.00	1530000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000004	CRS-004	Khóa học Ngữ Văn Tăng Cường 4	Nội dung chuẩn kiến thức Ngữ Văn	Khối 12	Ngữ Văn	24	90	36.00	1540000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000005	CRS-005	Khóa học Toán Tăng Cường 5	Nội dung chuẩn kiến thức Toán	Khối 8	Toán	24	90	36.00	1550000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000006	CRS-006	Khóa học Tiếng Anh Tăng Cường 6	Nội dung chuẩn kiến thức Tiếng Anh	Khối 9	Tiếng Anh	24	90	36.00	1560000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000007	CRS-007	Khóa học Vật Lý Tăng Cường 7	Nội dung chuẩn kiến thức Vật Lý	Khối 10	Vật Lý	24	90	36.00	1570000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000008	CRS-008	Khóa học Hóa Học Tăng Cường 8	Nội dung chuẩn kiến thức Hóa Học	Khối 11	Hóa Học	24	90	36.00	1580000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000009	CRS-009	Khóa học Ngữ Văn Tăng Cường 9	Nội dung chuẩn kiến thức Ngữ Văn	Khối 12	Ngữ Văn	24	90	36.00	1590000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000010	CRS-010	Khóa học Toán Tăng Cường 10	Nội dung chuẩn kiến thức Toán	Khối 8	Toán	24	90	36.00	1600000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000011	CRS-011	Khóa học Tiếng Anh Tăng Cường 11	Nội dung chuẩn kiến thức Tiếng Anh	Khối 9	Tiếng Anh	24	90	36.00	1610000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000012	CRS-012	Khóa học Vật Lý Tăng Cường 12	Nội dung chuẩn kiến thức Vật Lý	Khối 10	Vật Lý	24	90	36.00	1620000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000013	CRS-013	Khóa học Hóa Học Tăng Cường 13	Nội dung chuẩn kiến thức Hóa Học	Khối 11	Hóa Học	24	90	36.00	1630000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000014	CRS-014	Khóa học Ngữ Văn Tăng Cường 14	Nội dung chuẩn kiến thức Ngữ Văn	Khối 12	Ngữ Văn	24	90	36.00	1640000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000015	CRS-015	Khóa học Toán Tăng Cường 15	Nội dung chuẩn kiến thức Toán	Khối 8	Toán	24	90	36.00	1650000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000016	CRS-016	Khóa học Tiếng Anh Tăng Cường 16	Nội dung chuẩn kiến thức Tiếng Anh	Khối 9	Tiếng Anh	24	90	36.00	1660000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000017	CRS-017	Khóa học Vật Lý Tăng Cường 17	Nội dung chuẩn kiến thức Vật Lý	Khối 10	Vật Lý	24	90	36.00	1670000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000018	CRS-018	Khóa học Hóa Học Tăng Cường 18	Nội dung chuẩn kiến thức Hóa Học	Khối 11	Hóa Học	24	90	36.00	1680000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000019	CRS-019	Khóa học Ngữ Văn Tăng Cường 19	Nội dung chuẩn kiến thức Ngữ Văn	Khối 12	Ngữ Văn	24	90	36.00	1690000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000020	CRS-020	Khóa học Toán Tăng Cường 20	Nội dung chuẩn kiến thức Toán	Khối 8	Toán	24	90	36.00	1700000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000021	CRS-021	Khóa học Tiếng Anh Tăng Cường 21	Nội dung chuẩn kiến thức Tiếng Anh	Khối 9	Tiếng Anh	24	90	36.00	1710000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000022	CRS-022	Khóa học Vật Lý Tăng Cường 22	Nội dung chuẩn kiến thức Vật Lý	Khối 10	Vật Lý	24	90	36.00	1720000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000023	CRS-023	Khóa học Hóa Học Tăng Cường 23	Nội dung chuẩn kiến thức Hóa Học	Khối 11	Hóa Học	24	90	36.00	1730000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000024	CRS-024	Khóa học Ngữ Văn Tăng Cường 24	Nội dung chuẩn kiến thức Ngữ Văn	Khối 12	Ngữ Văn	24	90	36.00	1740000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
22222222-2222-2222-2222-000000000025	CRS-025	Khóa học Toán Tăng Cường 25	Nội dung chuẩn kiến thức Toán	Khối 8	Toán	24	90	36.00	1750000.00	ACTIVE	2026-02-26 03:01:23.806461+00	\N	\N
\.


--
-- Data for Name: enrollments; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.enrollments (id, class_id, student_id, status, approved_at, rejected_at, created_at, updated_at) FROM stdin;
ba68feca-6a7f-46e4-863e-af4b63aba9cc	66666666-6666-6666-6666-000000000002	55555555-5555-5555-5555-000000000001	ENROLLED	\N	\N	2026-03-06 08:38:07.83279+00	2026-03-06 08:38:07.834034+00
bc9d8fca-240d-44f7-9f0d-a2bd06eb70e1	66666666-6666-6666-6666-000000000002	55555555-5555-5555-5555-000000000002	ENROLLED	\N	\N	2026-03-06 08:38:07.856885+00	2026-03-06 08:38:07.857579+00
b47c0fe3-ad45-4b1c-8fc3-e0646e470911	66666666-6666-6666-6666-000000000002	55555555-5555-5555-5555-000000000003	ENROLLED	\N	\N	2026-03-06 08:38:07.858622+00	2026-03-06 08:38:07.859361+00
e2724999-b55f-450b-bb7b-7a9abe9c58ff	66666666-6666-6666-6666-000000000002	55555555-5555-5555-5555-000000000005	ENROLLED	\N	\N	2026-03-31 02:10:21.479662+00	2026-03-31 02:10:21.481514+00
\.


--
-- Data for Name: leave_requests; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.leave_requests (id, student_id, leave_type, apply_date, late_minutes, early_minutes, reason, documents, class_id, lesson_id, subject, status, approved_by_id, approved_at, rejection_reason, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: lesson_summaries; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.lesson_summaries (id, lesson_id, topic, lesson_content, class_feedback, homework, homework_deadline, teacher_notes, created_by_id, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: lessons; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.lessons (id, class_id, date_start, date_end, room_id, notes, created_at, updated_at, teacher_id) FROM stdin;
\.


--
-- Data for Name: objectives; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.objectives (id, code, name, program_id) FROM stdin;
\.


--
-- Data for Name: outcomes; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.outcomes (id, code, name, program_id, objective_id) FROM stdin;
\.


--
-- Data for Name: password_resets; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.password_resets (id, user_id, token_hash, expires_at, used_at, requested_ip, user_agent, created_at, updated_at, deleted_at) FROM stdin;
4405767b-6a43-47ee-916b-180cef7f91db	5b719ad5-b8c0-436b-b9dd-6bde07906dba	$2a$12$HLul7nu5hBw8gMFMEkto2ec6elnqk3LTLRFP3rVwEo0w4vDe.rSIO	2026-02-06 09:12:50.416942+00	\N			2026-02-06 08:57:50.417168+00	2026-02-06 08:57:50.417168+00	\N
9a8dc3f2-4ac2-408e-9bea-59dd4ccc04ee	5b719ad5-b8c0-436b-b9dd-6bde07906dba	$2a$12$6TTFCBdqyEDRwHGFtL.NZ.73ZAD5yQtDe4e6.NkDU0O1WwdvYvLfi	2026-02-06 09:15:09.132223+00	\N			2026-02-06 09:00:09.131451+00	2026-02-06 09:00:09.131451+00	\N
5aa01277-a950-4480-a108-bb06e4091b0a	5b719ad5-b8c0-436b-b9dd-6bde07906dba	$2a$12$IUmyiOCwBm5BYJTSgpiKmuNQTHVOXUejCCbUrQTgl1CasvpXwOGTi	2026-02-06 09:21:10.141002+00	\N			2026-02-06 09:06:10.141195+00	2026-02-06 09:06:10.141195+00	\N
\.


--
-- Data for Name: program_courses; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.program_courses (id, program_id, course_id) FROM stdin;
0a161bc6-0bb6-4f0b-8092-9c9c7504c3a5	33333333-3333-3333-3333-000000000001	22222222-2222-2222-2222-000000000001
2262ce0b-5525-4305-bebe-3839b7340a1f	33333333-3333-3333-3333-000000000001	22222222-2222-2222-2222-000000000002
46a51457-8695-4bed-a6a5-2823c2389574	33333333-3333-3333-3333-000000000001	22222222-2222-2222-2222-000000000003
508d0cb2-f37c-4294-a266-fea9516a1887	33333333-3333-3333-3333-000000000001	22222222-2222-2222-2222-000000000004
9b19fd85-2eeb-4420-aea2-4ccfc7c50fb0	33333333-3333-3333-3333-000000000001	22222222-2222-2222-2222-000000000005
9699ec75-cd04-41aa-9280-e0a52b0a4e50	33333333-3333-3333-3333-000000000001	22222222-2222-2222-2222-000000000008
\.


--
-- Data for Name: programs; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.programs (id, code, name, track, effective_from, effective_to, created_by_id, approved_by_id, approval_note, published_at, archived_at, created_at, updated_at, deleted_at, status) FROM stdin;
33333333-3333-3333-3333-000000000001	PRG-001	Chương trình Đào tạo 1 (ADVANCED)	ADVANCED	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000002	PRG-002	Chương trình Đào tạo 2 (SUPPORT)	SUPPORT	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000003	PRG-003	Chương trình Đào tạo 3 (BASIC)	BASIC	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000004	PRG-004	Chương trình Đào tạo 4 (ADVANCED)	ADVANCED	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000005	PRG-005	Chương trình Đào tạo 5 (SUPPORT)	SUPPORT	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000006	PRG-006	Chương trình Đào tạo 6 (BASIC)	BASIC	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000007	PRG-007	Chương trình Đào tạo 7 (ADVANCED)	ADVANCED	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000008	PRG-008	Chương trình Đào tạo 8 (SUPPORT)	SUPPORT	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000009	PRG-009	Chương trình Đào tạo 9 (BASIC)	BASIC	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000010	PRG-010	Chương trình Đào tạo 10 (ADVANCED)	ADVANCED	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000011	PRG-011	Chương trình Đào tạo 11 (SUPPORT)	SUPPORT	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000012	PRG-012	Chương trình Đào tạo 12 (BASIC)	BASIC	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000013	PRG-013	Chương trình Đào tạo 13 (ADVANCED)	ADVANCED	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000014	PRG-014	Chương trình Đào tạo 14 (SUPPORT)	SUPPORT	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000015	PRG-015	Chương trình Đào tạo 15 (BASIC)	BASIC	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000016	PRG-016	Chương trình Đào tạo 16 (ADVANCED)	ADVANCED	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000017	PRG-017	Chương trình Đào tạo 17 (SUPPORT)	SUPPORT	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000018	PRG-018	Chương trình Đào tạo 18 (BASIC)	BASIC	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000019	PRG-019	Chương trình Đào tạo 19 (ADVANCED)	ADVANCED	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000020	PRG-020	Chương trình Đào tạo 20 (SUPPORT)	SUPPORT	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000021	PRG-021	Chương trình Đào tạo 21 (BASIC)	BASIC	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000022	PRG-022	Chương trình Đào tạo 22 (ADVANCED)	ADVANCED	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000023	PRG-023	Chương trình Đào tạo 23 (SUPPORT)	SUPPORT	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000024	PRG-024	Chương trình Đào tạo 24 (BASIC)	BASIC	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
33333333-3333-3333-3333-000000000025	PRG-025	Chương trình Đào tạo 25 (ADVANCED)	ADVANCED	\N	\N	00000000-0000-0000-0000-000000000001	\N	\N	\N	\N	2026-02-26 03:06:23.606523+00	\N	\N	ACTIVE
\.


--
-- Data for Name: rooms; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.rooms (id, code, name, capacity, address, created_at, updated_at, deleted_at) FROM stdin;
dba1fc6a-9222-481e-8f28-e77df9966fcd		Online Room 1	30		2026-02-26 02:55:50.665376+00	2026-02-26 02:55:50.667759+00	\N
11111111-1111-1111-1111-000000000004	P104	Phòng Học 104	30	Tầng 1 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000005	P105	Phòng Học 105	30	Tầng 1 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000006	P106	Phòng Học 106	30	Tầng 1 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000007	P107	Phòng Học 107	30	Tầng 1 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000008	P108	Phòng Học 108	30	Tầng 1 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000009	P109	Phòng Học 109	30	Tầng 1 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000010	P110	Phòng Học 110	30	Tầng 1 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000011	P111	Phòng Học 111	30	Tầng 1 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000012	P112	Phòng Học 112	30	Tầng 1 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000013	P113	Phòng Học 113	30	Tầng 1 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000014	P114	Phòng Học 114	30	Tầng 1 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000015	P115	Phòng Học 115	30	Tầng 2 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000016	P116	Phòng Học 116	30	Tầng 2 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000017	P117	Phòng Học 117	30	Tầng 2 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000018	P118	Phòng Học 118	30	Tầng 2 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000019	P119	Phòng Học 119	30	Tầng 2 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000020	P120	Phòng Học 120	30	Tầng 2 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000021	P121	Phòng Học 121	30	Tầng 2 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000022	P122	Phòng Học 122	30	Tầng 2 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000023	P123	Phòng Học 123	30	Tầng 2 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000024	P124	Phòng Học 124	30	Tầng 2 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000025	P125	Phòng Học 125	30	Tầng 2 Tòa A	2026-02-26 03:01:23.776732+00	\N	\N
11111111-1111-1111-1111-000000000001	P101	Phòng Học 101	30	Tầng 1 Tòa A	2026-02-26 03:01:23.776732+00	2026-03-06 09:24:07.892881+00	\N
11111111-1111-1111-1111-000000000002	P102	Phòng Học 102	30	Tầng 1 Tòa A	2026-02-26 03:01:23.776732+00	2026-03-31 02:12:02.878977+00	\N
11111111-1111-1111-1111-000000000003	P103	Phòng Học 103	30	Tầng 1 Tòa A	2026-02-26 03:01:23.776732+00	2026-04-06 02:34:27.790199+00	\N
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.schema_migrations (version, dirty) FROM stdin;
19	f
\.


--
-- Data for Name: shifts; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.shifts (id, code, name, start_time, end_time, duration_minutes, session_type, is_active, notes, created_at, updated_at, deleted_at) FROM stdin;
c692a236-2b83-4846-8a3b-63a9bcb9dfb1	ca1	Ca số 1	07:00	09:00	120	MORNING	t	Ca học số 1, ca đầu tiên trong ngày	2026-04-14 03:57:22.635864+00	2026-04-14 03:57:22.636545+00	\N
\.


--
-- Data for Name: students; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.students (id, code, full_name, email, phone, guardian_phone, grade_level, status, date_of_birth, gender, address, created_at, updated_at, deleted_at) FROM stdin;
55555555-5555-5555-5555-000000000001	HS-001	HS Phạm Thu C	hs1@student.com	0800000001	0900000011	Khối 9	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000003	HS-003	HS Huỳnh Hải E	hs3@student.com	0800000003	0900000013	Khối 11	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000004	HS-004	HS Phan Đức F	hs4@student.com	0800000004	0900000014	Khối 12	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000005	HS-005	HS Vũ Bảo G	hs5@student.com	0800000005	0900000015	Khối 8	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000006	HS-006	HS Võ Lan H	hs6@student.com	0800000006	0900000016	Khối 9	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000007	HS-007	HS Đặng Văn I	hs7@student.com	0800000007	0900000017	Khối 10	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000008	HS-008	HS Nguyễn Thị K	hs8@student.com	0800000008	0900000018	Khối 11	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000009	HS-009	HS Trần Thanh Long	hs9@student.com	0800000009	0900000019	Khối 12	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000010	HS-010	HS Lê Minh Linh	hs10@student.com	0800000010	0900000020	Khối 8	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000011	HS-011	HS Phạm Thu Hùng	hs11@student.com	0800000011	0900000021	Khối 9	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000012	HS-012	HS Hoàng Ngọc Hường	hs12@student.com	0800000012	0900000022	Khối 10	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000013	HS-013	HS Huỳnh Hải Tùng	hs13@student.com	0800000013	0900000023	Khối 11	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000014	HS-014	HS Phan Đức A	hs14@student.com	0800000014	0900000024	Khối 12	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000015	HS-015	HS Vũ Bảo B	hs15@student.com	0800000015	0900000025	Khối 8	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000016	HS-016	HS Võ Lan C	hs16@student.com	0800000016	0900000026	Khối 9	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000017	HS-017	HS Đặng Văn D	hs17@student.com	0800000017	0900000027	Khối 10	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000018	HS-018	HS Nguyễn Thị E	hs18@student.com	0800000018	0900000028	Khối 11	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000019	HS-019	HS Trần Thanh F	hs19@student.com	0800000019	0900000029	Khối 12	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000020	HS-020	HS Lê Minh G	hs20@student.com	0800000020	0900000030	Khối 8	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000021	HS-021	HS Phạm Thu H	hs21@student.com	0800000021	0900000031	Khối 9	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000022	HS-022	HS Hoàng Ngọc I	hs22@student.com	0800000022	0900000032	Khối 10	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000023	HS-023	HS Huỳnh Hải K	hs23@student.com	0800000023	0900000033	Khối 11	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000024	HS-024	HS Phan Đức Long	hs24@student.com	0800000024	0900000034	Khối 12	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000025	HS-025	HS Vũ Bảo Linh	hs25@student.com	0800000025	0900000035	Khối 8	ACTIVE	\N	\N	\N	2026-02-26 03:06:23.647713+00	\N	\N
55555555-5555-5555-5555-000000000002	HS-002	HS Hoàng Ngọc D	hs2@student.com	0800000002	0900000012	Khối 10	ACTIVE	\N	MALE		2026-02-26 03:06:23.647713+00	2026-03-06 08:36:57.395646+00	\N
208c612b-513f-4489-a080-348d36d621b3	student1	Nguyen Van Student1	student1@gmail.com	0998267321	0998267321	10	ACTIVE	2010-02-25 00:00:00+00	MALE	Hoan Kiem, Ha Noi	2026-02-25 03:39:21.070614+00	2026-02-26 02:57:56.52406+00	2026-03-31 02:04:53.498042+00
\.


--
-- Data for Name: teachers; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.teachers (id, code, full_name, email, phone, is_school_teacher, school_name, employment_type, status, notes, created_at, updated_at, deleted_at) FROM stdin;
1c336210-2aee-4d4a-a322-801b675f9340	teacher11	Nguyen Van A	nguyenvana@gmail.com	012345678	t	Trung Tam Tieng Anh	PART_TIME	ACTIVE	trung tam tieng anh	2026-02-09 08:01:31.578526+00	2026-02-09 08:01:31.57909+00	\N
bb7c8104-f937-4472-9ef7-64cb3458a4ec	teacher12	Nguyen Thi B	nguyenthib@gmail.com	098776543	t	ABC	FULL_TIME	ACTIVE	ABC	2026-02-09 08:16:37.430742+00	2026-02-09 08:16:37.431708+00	\N
44444444-4444-4444-4444-000000000001	GV-001	GV Trần Thị B	gv1@educenter.com	0900000001	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000002	GV-002	GV Lê Thanh C	gv2@educenter.com	0900000002	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000003	GV-003	GV Phạm Minh D	gv3@educenter.com	0900000003	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000004	GV-004	GV Hoàng Thu E	gv4@educenter.com	0900000004	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000005	GV-005	GV Huỳnh Ngọc F	gv5@educenter.com	0900000005	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000006	GV-006	GV Phan Hải G	gv6@educenter.com	0900000006	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000007	GV-007	GV Vũ Đức H	gv7@educenter.com	0900000007	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000008	GV-008	GV Võ Bảo I	gv8@educenter.com	0900000008	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000009	GV-009	GV Đặng Lan K	gv9@educenter.com	0900000009	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000010	GV-010	GV Nguyễn Văn Long	gv10@educenter.com	0900000010	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000011	GV-011	GV Trần Thị Linh	gv11@educenter.com	0900000011	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000012	GV-012	GV Lê Thanh Hùng	gv12@educenter.com	0900000012	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000013	GV-013	GV Phạm Minh Hường	gv13@educenter.com	0900000013	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000014	GV-014	GV Hoàng Thu Tùng	gv14@educenter.com	0900000014	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000015	GV-015	GV Huỳnh Ngọc A	gv15@educenter.com	0900000015	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000016	GV-016	GV Phan Hải B	gv16@educenter.com	0900000016	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000017	GV-017	GV Vũ Đức C	gv17@educenter.com	0900000017	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000018	GV-018	GV Võ Bảo D	gv18@educenter.com	0900000018	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000019	GV-019	GV Đặng Lan E	gv19@educenter.com	0900000019	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000020	GV-020	GV Nguyễn Văn F	gv20@educenter.com	0900000020	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000021	GV-021	GV Trần Thị G	gv21@educenter.com	0900000021	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000022	GV-022	GV Lê Thanh H	gv22@educenter.com	0900000022	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000023	GV-023	GV Phạm Minh I	gv23@educenter.com	0900000023	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000024	GV-024	GV Hoàng Thu K	gv24@educenter.com	0900000024	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
44444444-4444-4444-4444-000000000025	GV-025	GV Huỳnh Ngọc Long	gv25@educenter.com	0900000025	f	\N	FULL_TIME	ACTIVE	\N	2026-02-26 03:06:23.629184+00	\N	\N
e477249f-de22-4876-8e05-c22b5f95ec3c	teacherA	Nguyen Van A	nguyenvanaaa@gmail.com	076867856567	t	Trung Tam Tieng Anh AAA	FULL_TIME	ACTIVE	Trung Tam Tieng Anh AAA	2026-03-31 02:23:26.984959+00	2026-03-31 02:23:26.98605+00	\N
\.


--
-- Data for Name: user_otps; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.user_otps (id, user_id, otp_hash, expired_at, used_at, created_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.users (id, code, full_name, email, password, role, is_active, created_at, updated_at, deleted_at) FROM stdin;
36956f7a-936f-457a-bfcb-b0dc8539b96e	admin	admin	admin	$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R	ADMIN	t	2026-01-27 10:31:55.969627+00	2026-01-27 10:31:55.970449+00	\N
664c9301-98d5-4c2c-9518-5e12f8da3d66	user1	user1	user1	$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R	STUDENT	t	2026-01-27 10:31:55.969627+00	2026-01-27 10:31:55.970449+00	\N
2c2dcf01-14be-4f4b-b9ad-d6a978c11cdb	user2	user2	user2	$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R	STUDENT	t	2026-01-27 10:31:55.969627+00	2026-01-27 10:31:55.970449+00	\N
46adda73-bfb6-4200-8955-01de43c048a6	testuser	testuser	testuser	$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R	STUDENT	t	2026-01-27 10:31:55.969627+00	2026-01-27 10:31:55.970449+00	\N
477a45a8-0b65-4d5a-a45e-fc2831e2de2c	john_doe	john_doe	john_doe	$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R	STUDENT	t	2026-01-27 10:31:55.969627+00	2026-01-27 10:31:55.970449+00	\N
0a4662ed-fabf-4959-a57a-ecdaef57542a	jane_smith	jane_smith	jane_smith	$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R	STUDENT	t	2026-01-27 10:31:55.969627+00	2026-01-27 10:31:55.970449+00	\N
3417fedc-4b3a-43f1-8840-48302fe4fbc9	bob_wilson	bob_wilson	bob_wilson	$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R	STUDENT	t	2026-01-27 10:31:55.969627+00	2026-01-27 10:31:55.970449+00	\N
0b17a289-ccdd-47af-a868-6f4fa5a8a918	alice_jones	alice_jones	alice_jones	$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R	STUDENT	t	2026-01-27 10:31:55.969627+00	2026-01-27 10:31:55.970449+00	\N
1deb366e-8dfe-4d0e-9742-d0c79629a5a1	charlie_brown	charlie_brown	charlie_brown	$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R	STUDENT	t	2026-01-27 10:31:55.969627+00	2026-01-27 10:31:55.970449+00	\N
3015df51-120c-4031-a6b3-32972d470893	david_miller	david_miller	david_miller	$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R	STUDENT	t	2026-01-27 10:31:55.969627+00	2026-01-27 10:31:55.970449+00	\N
d6fa5960-dc66-4131-b07e-8ceb6b5a5b23	TEACHER001	John Teacher	teacher@example.com	$2a$10$95X/d6EwOGozE5lTYAxBD.g.TGD3i/ptPGN1xR9C7UIEeE8TXRkzG	TEACHER	t	2026-01-27 10:32:43.8088+00	2026-01-27 10:32:43.8088+00	\N
41ba9a7f-05ef-4ef1-a5a9-7bdad8537fe6	STUDENT001	Alice Student	student1@example.com	$2a$10$95X/d6EwOGozE5lTYAxBD.g.TGD3i/ptPGN1xR9C7UIEeE8TXRkzG	STUDENT	t	2026-01-27 10:32:43.808801+00	2026-01-27 10:32:43.808801+00	\N
5872d694-4128-4ba8-9cf1-231db4f6946d	STUDENT002	Bob Student	student2@example.com	$2a$10$95X/d6EwOGozE5lTYAxBD.g.TGD3i/ptPGN1xR9C7UIEeE8TXRkzG	STUDENT	t	2026-01-27 10:32:43.808802+00	2026-01-27 10:32:43.808802+00	\N
40d5cbeb-5274-494d-9597-820e021b4e1c	STUDENT003	Charlie Student	student3@example.com	$2a$10$95X/d6EwOGozE5lTYAxBD.g.TGD3i/ptPGN1xR9C7UIEeE8TXRkzG	STUDENT	t	2026-01-27 10:32:43.808802+00	2026-01-27 10:32:43.808803+00	\N
5b719ad5-b8c0-436b-b9dd-6bde07906dba	dev	Nguyen The Ha	thehanguyen02@gmail.com	$2a$12$rH.rj.Tf1IftjeOJAAOP7uUJy9vN//vMkcAUBqk6ri868ZQDoJFs.	ADMIN	t	2026-01-28 09:50:22.05681+00	2026-01-28 09:50:22.058218+00	\N
3496c9ea-6060-49c3-82aa-bdc04bff6cf6	ADMIN001	Admin User	admin@example.com	$2a$10$95X/d6EwOGozE5lTYAxBD.g.TGD3i/ptPGN1xR9C7UIEeE8TXRkzG	ADMIN	t	2026-01-27 10:32:43.808799+00	2026-01-27 10:32:43.808799+00	\N
00000000-0000-0000-0000-000000000001	ADM-01	Quản trị viên	admin@educenter.com	hashed_pw	ADMIN	t	2026-02-26 03:01:23.747398+00	\N	\N
\.


--
-- Name: academic_records academic_records_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.academic_records
    ADD CONSTRAINT academic_records_pkey PRIMARY KEY (id);


--
-- Name: attendances attendances_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.attendances
    ADD CONSTRAINT attendances_pkey PRIMARY KEY (id);


--
-- Name: class_schedules class_schedules_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.class_schedules
    ADD CONSTRAINT class_schedules_pkey PRIMARY KEY (id);


--
-- Name: classes classes_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.classes
    ADD CONSTRAINT classes_pkey PRIMARY KEY (id);


--
-- Name: consultations consultations_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.consultations
    ADD CONSTRAINT consultations_pkey PRIMARY KEY (id);


--
-- Name: courses courses_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.courses
    ADD CONSTRAINT courses_pkey PRIMARY KEY (id);


--
-- Name: enrollments enrollments_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.enrollments
    ADD CONSTRAINT enrollments_pkey PRIMARY KEY (id);


--
-- Name: leave_requests leave_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.leave_requests
    ADD CONSTRAINT leave_requests_pkey PRIMARY KEY (id);


--
-- Name: lesson_summaries lesson_summaries_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.lesson_summaries
    ADD CONSTRAINT lesson_summaries_pkey PRIMARY KEY (id);


--
-- Name: lessons lessons_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.lessons
    ADD CONSTRAINT lessons_pkey PRIMARY KEY (id);


--
-- Name: objectives objectives_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.objectives
    ADD CONSTRAINT objectives_pkey PRIMARY KEY (id);


--
-- Name: outcomes outcomes_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.outcomes
    ADD CONSTRAINT outcomes_pkey PRIMARY KEY (id);


--
-- Name: password_resets password_resets_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.password_resets
    ADD CONSTRAINT password_resets_pkey PRIMARY KEY (id);


--
-- Name: program_courses program_courses_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.program_courses
    ADD CONSTRAINT program_courses_pkey PRIMARY KEY (id);


--
-- Name: programs programs_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.programs
    ADD CONSTRAINT programs_pkey PRIMARY KEY (id);


--
-- Name: rooms rooms_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.rooms
    ADD CONSTRAINT rooms_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: shifts shifts_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.shifts
    ADD CONSTRAINT shifts_pkey PRIMARY KEY (id);


--
-- Name: students students_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.students
    ADD CONSTRAINT students_pkey PRIMARY KEY (id);


--
-- Name: teachers teachers_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.teachers
    ADD CONSTRAINT teachers_pkey PRIMARY KEY (id);


--
-- Name: classes uni_classes_code; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.classes
    ADD CONSTRAINT uni_classes_code UNIQUE (code);


--
-- Name: courses uni_courses_code; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.courses
    ADD CONSTRAINT uni_courses_code UNIQUE (code);


--
-- Name: lesson_summaries uni_lesson_summaries_lesson_id; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.lesson_summaries
    ADD CONSTRAINT uni_lesson_summaries_lesson_id UNIQUE (lesson_id);


--
-- Name: objectives uni_objectives_code; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.objectives
    ADD CONSTRAINT uni_objectives_code UNIQUE (code);


--
-- Name: outcomes uni_outcomes_code; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.outcomes
    ADD CONSTRAINT uni_outcomes_code UNIQUE (code);


--
-- Name: programs uni_programs_code; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.programs
    ADD CONSTRAINT uni_programs_code UNIQUE (code);


--
-- Name: rooms uni_rooms_code; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.rooms
    ADD CONSTRAINT uni_rooms_code UNIQUE (code);


--
-- Name: students uni_students_code; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.students
    ADD CONSTRAINT uni_students_code UNIQUE (code);


--
-- Name: teachers uni_teachers_code; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.teachers
    ADD CONSTRAINT uni_teachers_code UNIQUE (code);


--
-- Name: users uni_users_code; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_code UNIQUE (code);


--
-- Name: users uni_users_email; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_email UNIQUE (email);


--
-- Name: user_otps user_otps_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.user_otps
    ADD CONSTRAINT user_otps_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_classes_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_classes_deleted_at ON public.classes USING btree (deleted_at);


--
-- Name: idx_courses_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_courses_deleted_at ON public.courses USING btree (deleted_at);


--
-- Name: idx_password_resets_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_password_resets_deleted_at ON public.password_resets USING btree (deleted_at);


--
-- Name: idx_programs_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_programs_deleted_at ON public.programs USING btree (deleted_at);


--
-- Name: idx_shifts_code; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_shifts_code ON public.shifts USING btree (code);


--
-- Name: idx_shifts_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_shifts_deleted_at ON public.shifts USING btree (deleted_at);


--
-- Name: idx_students_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_students_deleted_at ON public.students USING btree (deleted_at);


--
-- Name: idx_teachers_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_teachers_deleted_at ON public.teachers USING btree (deleted_at);


--
-- Name: idx_user_otps_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_user_otps_deleted_at ON public.user_otps USING btree (deleted_at);


--
-- Name: idx_user_otps_user_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_user_otps_user_id ON public.user_otps USING btree (user_id);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: academic_records fk_academic_records_lesson_summary; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.academic_records
    ADD CONSTRAINT fk_academic_records_lesson_summary FOREIGN KEY (lesson_summary_id) REFERENCES public.lesson_summaries(id) ON DELETE CASCADE;


--
-- Name: academic_records fk_academic_records_student; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.academic_records
    ADD CONSTRAINT fk_academic_records_student FOREIGN KEY (student_id) REFERENCES public.students(id) ON DELETE CASCADE;


--
-- Name: attendances fk_attendances_lesson; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.attendances
    ADD CONSTRAINT fk_attendances_lesson FOREIGN KEY (lesson_id) REFERENCES public.lessons(id) ON DELETE CASCADE;


--
-- Name: attendances fk_attendances_student; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.attendances
    ADD CONSTRAINT fk_attendances_student FOREIGN KEY (student_id) REFERENCES public.students(id) ON DELETE CASCADE;


--
-- Name: class_schedules fk_class_schedules_class; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.class_schedules
    ADD CONSTRAINT fk_class_schedules_class FOREIGN KEY (class_id) REFERENCES public.classes(id) ON DELETE CASCADE;


--
-- Name: class_schedules fk_class_schedules_room; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.class_schedules
    ADD CONSTRAINT fk_class_schedules_room FOREIGN KEY (room_id) REFERENCES public.rooms(id);


--
-- Name: class_schedules fk_class_schedules_shift; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.class_schedules
    ADD CONSTRAINT fk_class_schedules_shift FOREIGN KEY (shift_id) REFERENCES public.shifts(id);


--
-- Name: class_schedules fk_classes_class_schedules; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.class_schedules
    ADD CONSTRAINT fk_classes_class_schedules FOREIGN KEY (class_id) REFERENCES public.classes(id);


--
-- Name: classes fk_classes_course; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.classes
    ADD CONSTRAINT fk_classes_course FOREIGN KEY (course_id) REFERENCES public.courses(id);


--
-- Name: classes fk_classes_program; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.classes
    ADD CONSTRAINT fk_classes_program FOREIGN KEY (program_id) REFERENCES public.programs(id);


--
-- Name: classes fk_classes_room; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.classes
    ADD CONSTRAINT fk_classes_room FOREIGN KEY (room_id) REFERENCES public.rooms(id);


--
-- Name: classes fk_classes_teacher; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.classes
    ADD CONSTRAINT fk_classes_teacher FOREIGN KEY (teacher_id) REFERENCES public.teachers(id);


--
-- Name: enrollments fk_enrollments_class; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.enrollments
    ADD CONSTRAINT fk_enrollments_class FOREIGN KEY (class_id) REFERENCES public.classes(id) ON DELETE CASCADE;


--
-- Name: enrollments fk_enrollments_student; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.enrollments
    ADD CONSTRAINT fk_enrollments_student FOREIGN KEY (student_id) REFERENCES public.students(id) ON DELETE CASCADE;


--
-- Name: leave_requests fk_leave_requests_approved_by; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.leave_requests
    ADD CONSTRAINT fk_leave_requests_approved_by FOREIGN KEY (approved_by_id) REFERENCES public.users(id);


--
-- Name: leave_requests fk_leave_requests_class; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.leave_requests
    ADD CONSTRAINT fk_leave_requests_class FOREIGN KEY (class_id) REFERENCES public.classes(id) ON DELETE SET NULL;


--
-- Name: leave_requests fk_leave_requests_lesson; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.leave_requests
    ADD CONSTRAINT fk_leave_requests_lesson FOREIGN KEY (lesson_id) REFERENCES public.lessons(id) ON DELETE SET NULL;


--
-- Name: leave_requests fk_leave_requests_student; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.leave_requests
    ADD CONSTRAINT fk_leave_requests_student FOREIGN KEY (student_id) REFERENCES public.students(id) ON DELETE CASCADE;


--
-- Name: lesson_summaries fk_lesson_summaries_created_by; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.lesson_summaries
    ADD CONSTRAINT fk_lesson_summaries_created_by FOREIGN KEY (created_by_id) REFERENCES public.users(id);


--
-- Name: lesson_summaries fk_lesson_summaries_lesson; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.lesson_summaries
    ADD CONSTRAINT fk_lesson_summaries_lesson FOREIGN KEY (lesson_id) REFERENCES public.lessons(id) ON DELETE CASCADE;


--
-- Name: lessons fk_lessons_class; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.lessons
    ADD CONSTRAINT fk_lessons_class FOREIGN KEY (class_id) REFERENCES public.classes(id) ON DELETE CASCADE;


--
-- Name: lessons fk_lessons_room; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.lessons
    ADD CONSTRAINT fk_lessons_room FOREIGN KEY (room_id) REFERENCES public.rooms(id);


--
-- Name: lessons fk_lessons_teacher; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.lessons
    ADD CONSTRAINT fk_lessons_teacher FOREIGN KEY (teacher_id) REFERENCES public.teachers(id);


--
-- Name: objectives fk_objectives_program; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.objectives
    ADD CONSTRAINT fk_objectives_program FOREIGN KEY (program_id) REFERENCES public.programs(id) ON DELETE CASCADE;


--
-- Name: outcomes fk_outcomes_objective; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.outcomes
    ADD CONSTRAINT fk_outcomes_objective FOREIGN KEY (objective_id) REFERENCES public.objectives(id) ON DELETE SET NULL;


--
-- Name: outcomes fk_outcomes_program; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.outcomes
    ADD CONSTRAINT fk_outcomes_program FOREIGN KEY (program_id) REFERENCES public.programs(id) ON DELETE CASCADE;


--
-- Name: program_courses fk_program_courses_course; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.program_courses
    ADD CONSTRAINT fk_program_courses_course FOREIGN KEY (course_id) REFERENCES public.courses(id) ON DELETE CASCADE;


--
-- Name: program_courses fk_program_courses_program; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.program_courses
    ADD CONSTRAINT fk_program_courses_program FOREIGN KEY (program_id) REFERENCES public.programs(id) ON DELETE CASCADE;


--
-- Name: programs fk_programs_approved_by; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.programs
    ADD CONSTRAINT fk_programs_approved_by FOREIGN KEY (approved_by_id) REFERENCES public.users(id);


--
-- Name: programs fk_programs_created_by; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.programs
    ADD CONSTRAINT fk_programs_created_by FOREIGN KEY (created_by_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

\unrestrict DdX62PsYk80fV80n3jObLeAKtAWXehdJI5e0N4kS5jtdnMd5yqMXkgP3cQOg3w4

