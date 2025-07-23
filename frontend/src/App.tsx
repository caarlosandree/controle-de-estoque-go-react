import { Routes, Route } from 'react-router-dom';
import { Layout } from './components/Layout';
import { ProductListPage } from './pages/ProductListPage';

function App() {
    return (
        <Routes>
            <Route path="/" element={<Layout />}>
                <Route index element={<ProductListPage />} />
            </Route>
        </Routes>
    )
}

export default App