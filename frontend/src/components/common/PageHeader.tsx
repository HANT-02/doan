import React from 'react';
import { Box, Typography, Stack, Breadcrumbs, Link } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import NavigateNextIcon from '@mui/icons-material/NavigateNext';

interface PageHeaderProps {
  title: string;
  subtitle?: string;
  actions?: React.ReactNode;
  breadcrumbs?: { label: string; path?: string }[];
}

const PageHeader: React.FC<PageHeaderProps> = ({ title, subtitle, actions, breadcrumbs }) => {
  return (
    <Box sx={{ mb: 4 }}>
      {breadcrumbs && (
        <Breadcrumbs
          separator={<NavigateNextIcon fontSize="small" />}
          sx={{ mb: 1, '& .MuiBreadcrumbs-li': { fontSize: '0.875rem' } }}
        >
          {breadcrumbs.map((bc, index) => (
            bc.path ? (
              <Link
                key={index}
                component={RouterLink}
                to={bc.path}
                underline="hover"
                color="inherit"
              >
                {bc.label}
              </Link>
            ) : (
              <Typography key={index} color="text.secondary" sx={{ fontSize: '0.875rem' }}>
                {bc.label}
              </Typography>
            )
          ))}
        </Breadcrumbs>
      )}

      <Stack
        direction={{ xs: 'column', sm: 'row' }}
        justifyContent="space-between"
        alignItems={{ xs: 'flex-start', sm: 'center' }}
        spacing={2}
      >
        <Box>
          <Typography variant="h4" component="h1" gutterBottom sx={{ fontWeight: 700 }}>
            {title}
          </Typography>
          {subtitle && (
            <Typography variant="body1" color="text.secondary">
              {subtitle}
            </Typography>
          )}
        </Box>
        {actions && (
          <Box sx={{ alignSelf: { xs: 'flex-end', sm: 'center' } }}>
            {actions}
          </Box>
        )}
      </Stack>
    </Box>
  );
};

export default PageHeader;
