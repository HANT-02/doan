import React from 'react';
import { Box, CircularProgress } from '@mui/material';

const FullScreenLoader: React.FC = () => {
  return (
    <Box
      sx={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
        width: '100vw',
        position: 'fixed',
        top: 0,
        left: 0,
        zIndex: (theme) => theme.zIndex.drawer + 2, // Above drawer/sidebar
        backgroundColor: 'rgba(255, 255, 255, 0.8)', // Semi-transparent white background
      }}
    >
      <CircularProgress size={60} />
    </Box>
  );
};

export default FullScreenLoader;
