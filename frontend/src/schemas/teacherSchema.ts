import { z } from 'zod';

export const createTeacherSchema = z.object({
    full_name: z.string().min(1, 'Full name is required').max(255, 'Full name is too long'),
    email: z.string().email('Invalid email format').optional().or(z.literal('')),
    phone: z.string().max(20, 'Phone number is too long').optional().or(z.literal('')),
    code: z.string().max(50, 'Code is too long').optional().or(z.literal('')),
    is_school_teacher: z.boolean().optional(),
    school_name: z.string().max(255, 'School name is too long').optional().or(z.literal('')),
    employment_type: z.enum(['PART_TIME', 'FULL_TIME']).optional(),
    status: z.enum(['ACTIVE', 'INACTIVE']).optional(),
    notes: z.string().optional().or(z.literal('')),
});

export const updateTeacherSchema = z.object({
    full_name: z.string().min(1, 'Full name is required').max(255, 'Full name is too long').optional(),
    email: z.string().email('Invalid email format').optional().or(z.literal('')),
    phone: z.string().max(20, 'Phone number is too long').optional().or(z.literal('')),
    code: z.string().max(50, 'Code is too long').optional().or(z.literal('')),
    is_school_teacher: z.boolean().optional(),
    school_name: z.string().max(255, 'School name is too long').optional().or(z.literal('')),
    employment_type: z.enum(['PART_TIME', 'FULL_TIME']).optional(),
    status: z.enum(['ACTIVE', 'INACTIVE']).optional(),
    notes: z.string().optional().or(z.literal('')),
});

export type CreateTeacherFormData = z.infer<typeof createTeacherSchema>;
export type UpdateTeacherFormData = z.infer<typeof updateTeacherSchema>;
