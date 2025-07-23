import axios from 'axios';

// Cria uma instância do axios com a URL base da nossa API Go.
// Todas as requisições feitas com esta instância irão para http://localhost:8080
const api = axios.create({
    baseURL: 'http://localhost:8080',
});

export default api;