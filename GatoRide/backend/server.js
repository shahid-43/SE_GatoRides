import express from 'express';
import connectDB from './config/db.js';
import router from './routes/authorization.js';
import cors from 'cors';
import dotenv from 'dotenv';

dotenv.config();

const app = express();

// Connect to MongoDB
connectDB();
//iCPpHajYSS6WNsGa
// Middleware
app.use(cors());
app.use(express.json());

// Routes
app.use('/api/auth', router);
app.get('/', (req, res) => {
    res.send('Server is running...');
});
const PORT = process.env.PORT || 2000;

app.listen(PORT, () => console.log(`Server running on port ${PORT}`),
console.log('http://localhost:2000'));