import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import DashboardPage from './pages/DashboardPage';
import NewAlgorithmPostPage from './pages/NewAlgorithmPostPage';
import SettingsPage from './pages/SettingsPage';

export default function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Navigate to="/posts" replace />} />
                <Route path="/posts" element={<DashboardPage />} />
                <Route path="/posts/new/algorithm" element={<NewAlgorithmPostPage />} />
                <Route path="/posts/:id/edit" element={<NewAlgorithmPostPage />} />
                <Route path="/settings" element={<SettingsPage />} />
            </Routes>
        </BrowserRouter>
    );
}
