import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import MonitorsPage from "./app/monitors/page";
import NewMonitor from "./app/monitors/new/page";
import SettingsPage from "./app/settings/page";
import { Routes, Route, Navigate } from "react-router-dom";
import { client } from "./api/client.gen";
import ProxiesPage from "./app/proxies/page";
import NewProxy from "./app/proxies/new/page";
import NotifiersPage from "./app/notifiers/page";
import NewNotifier from "./app/notifiers/new/page";
import MonitorPage from "./app/monitors/view/page";
import { ThemeProvider } from "@/components/theme-provider";
import EditMonitor from "./app/monitors/edit/page";
import SHLoginPage from "./app/sh/login/page";
import SHRegisterPage from "./app/sh/register/page";
import { useAuthStore } from "@/store/auth";
import { setupInterceptors } from "./interceptors";
import { WebSocketProvider } from "./context/WebsocketContext";
import StatusPagesPage from "./app/status-pages/page";
import NewStatusPage from "./app/status-pages/new/page";
import SecurityPage from "./app/security/page";
import EditNotifier from "./app/notifiers/edit/page";
import EditProxy from "./app/proxies/edit/page";
import { TimezoneProvider } from './context/TimezoneContext';
// import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

export const configureClient = () => {
  const accessToken = useAuthStore.getState().accessToken;

  client.setConfig({
    baseURL: import.meta.env.VITE_API_URL + "/api/v1",
    headers: accessToken
      ? {
          Authorization: `Bearer ${accessToken}`,
        }
      : undefined,
  });
};

configureClient();
setupInterceptors();

useAuthStore.subscribe((state) => {
  client.setConfig({
    baseURL: import.meta.env.VITE_API_URL + "/api/v1",
    headers: state.accessToken
      ? {
          Authorization: `Bearer ${state.accessToken}`,
        }
      : undefined,
  });
});

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      staleTime: 1000 * 60 * 5, // 5 minutes
    },
    mutations: {
      retry: false,
    },
  },
});

export default function App() {
  const accessToken = useAuthStore((state) => state.accessToken);
  // const isDev = import.meta.env.MODE === "development";

  return (
    <ThemeProvider defaultTheme="dark" storageKey="peekaping-ui-theme">
      <TimezoneProvider>
        <QueryClientProvider client={queryClient}>
          <WebSocketProvider>
            <Routes>
              {!accessToken ? (
                <>
                  <Route path="/login" element={<SHLoginPage />} />
                  <Route path="/register" element={<SHRegisterPage />} />
                  <Route path="*" element={<Navigate to="/login" replace />} />
                </>
              ) : (
                <>
                  <Route path="/monitors" element={<MonitorsPage />} />
                  <Route path="/monitors/:id" element={<MonitorPage />} />
                  <Route path="/monitors/new" element={<NewMonitor />} />
                  <Route path="/monitors/edit/:id" element={<EditMonitor />} />

                  <Route path="/status-pages" element={<StatusPagesPage />} />
                  <Route path="/status-pages/new" element={<NewStatusPage />} />

                  <Route path="/proxies" element={<ProxiesPage />} />
                  <Route path="/proxies/new" element={<NewProxy />} />
                  <Route path="/proxies/edit/:id" element={<EditProxy />} />

                  <Route path="/notifiers" element={<NotifiersPage />} />
                  <Route path="/notifiers/new" element={<NewNotifier />} />
                  <Route path="/notifiers/:id/edit" element={<EditNotifier />} />

                  <Route path="/settings" element={<SettingsPage />} />
                  <Route path="/security" element={<SecurityPage />} />

                  <Route path="*" element={<Navigate to="/monitors" replace />} />
                </>
              )}
            </Routes>
          </WebSocketProvider>
          {/* {isDev && <ReactQueryDevtools initialIsOpen={false} />} */}
        </QueryClientProvider>
      </TimezoneProvider>
    </ThemeProvider>
  );
}
