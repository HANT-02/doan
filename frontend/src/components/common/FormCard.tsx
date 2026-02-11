import React from 'react';
import { Paper, Box, Typography, Divider } from '@mui/material';

interface FormCardProps {
  title?: string;
  subtitle?: string;
  children: React.ReactNode;
  footer?: React.ReactNode;
  sx?: any;
}

const FormCard: React.FC<FormCardProps> = ({ title, subtitle, children, footer, sx }) => {
  return (
    <Paper sx={{ overflow: 'hidden', ...sx }}>
      {(title || subtitle) && (
        <Box sx={{ px: 3, py: 2 }}>
          {title && (
            <Typography variant="h6" component="h2" sx={{ fontWeight: 600 }}>
              {title}
            </Typography>
          )}
          {subtitle && (
            <Typography variant="body2" color="text.secondary">
              {subtitle}
            </Typography>
          )}
        </Box>
      )}
      {(title || subtitle) && <Divider />}
      <Box sx={{ p: 3 }}>
        {children}
      </Box>
      {footer && (
        <>
          <Divider />
          <Box sx={{ px: 3, py: 2, bgcolor: 'background.default' }}>
            {footer}
          </Box>
        </>
      )}
    </Paper>
  );
};

export default FormCard;
