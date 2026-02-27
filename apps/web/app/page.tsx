"use client";

import {
   Box,
   Container,
   Typography,
   Card,
   CardContent,
   Stack,
} from "@mui/material";
import ThemeToggle from "@/components/ui/ThemeToggle";
import { useThemeMode } from "@/components/providers/ThemeProvider";

const Home = () => {
   const { mode } = useThemeMode();

   return (
      <Container maxWidth="lg">
         <Box sx={{ py: 8 }}>
            <Stack
               direction="row"
               justifyContent="space-between"
               alignItems="center"
               sx={{ mb: 4 }}
            >
               <Typography variant="h2" component="h1">
                  ProvisionHub
               </Typography>
               <ThemeToggle />
            </Stack>

            <Card>
               <CardContent>
                  <Typography variant="h5" gutterBottom>
                     Welcome to ProvisionHub
                  </Typography>
                  <Typography variant="body1" color="text.secondary">
                     Self-service Platform Provisioning • Git-native • Async
                  </Typography>
                  <Typography variant="body2" sx={{ mt: 2 }}>
                     Current theme: <strong>{mode}</strong>
                  </Typography>
               </CardContent>
            </Card>
         </Box>
      </Container>
   );
};

export default Home;
