import React, { useState, useEffect, useCallback } from 'react'; // <-- 1. Tambahkan useCallback
import { 
    PieChart, Pie, Cell, ResponsiveContainer, 
    BarChart, Bar, XAxis, YAxis, Tooltip, Legend, CartesianGrid 
} from 'recharts';
import { Loader, AlertCircle, DollarSign, Target, Activity, TrendingUp, BarChart as BarIcon } from 'lucide-react';

// --- MOCK DATA LOKAL (DIHAPUS) ---
// const mockSummary = { ... };
// const mockCharts = { ... };
// --- AKHIR MOCK DATA LOKAL ---


// API Base URL (Sekarang akan kita gunakan)
const API_BASE_URL = 'http://localhost:8080/api/v1';

// --- KOMPONEN LOKAL (StatCard dan ChartCard tidak berubah) ---
const StatCard = ({ title, value, icon: Icon, colorClass = 'text-gray-900', bgColorClass = 'bg-gray-50' }) => (
    <div className={`p-5 rounded-xl shadow-lg border border-gray-200 ${bgColorClass}`}>
        <div className="flex items-center justify-between mb-2">
            <p className="text-sm font-medium text-gray-500">{title}</p>
            {Icon && <Icon className={`w-6 h-6 ${colorClass.replace('text-', 'text-')}`} />}
        </div>
        <p className={`text-3xl font-bold ${colorClass}`}>{value}</p>
    </div>
);

const ChartCard = ({ title, children, icon: Icon }) => (
  <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200 h-full flex flex-col">
    <div className="flex items-center justify-between border-b pb-3 mb-4">
        <h3 className="text-xl font-semibold text-gray-800">{title}</h3>
        {Icon && <Icon className="w-6 h-6 text-indigo-500" />}
    </div>
    <div className="flex-grow w-full h-72 min-h-[288px]"> 
      {children}
    </div>
  </div>
);
// --- AKHIR KOMPONEN LOKAL ---


/**
 * Komponen Halaman Dashboard
 */
const Dashboard = () => {
    
    // State (tidak berubah)
    const [summary, setSummary] = useState({
        total_belanja: 0,
        budget: 0,
        sisa_budget: 0,
    });
    const [pieData, setPieData] = useState([]);
    const [barData, setBarData] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);

    
    // --- PERBAIKAN 2: Tambahkan helper fetchWithAuth ---
    const fetchWithAuth = useCallback(async (url, options = {}) => {
        const token = localStorage.getItem('authToken');
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers,
        };
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
        const response = await fetch(url, { ...options, headers });
        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.message || errorData.error || `Error ${response.status}: ${response.statusText}`);
        }
        const text = await response.text();
        return text ? JSON.parse(text) : null;
    }, []);


    // --- PERBAIKAN 3: Ganti simulateFetchData dengan fetchData asli ---
    /**
     * Mengambil data dashboard dari API backend
     */
    const fetchData = useCallback(async () => {
        setIsLoading(true);
        setError(null);

        try {
            // Panggil kedua endpoint secara bersamaan
            const [summaryData, chartsData] = await Promise.all([
                fetchWithAuth(`${API_BASE_URL}/dashboard/summary`),
                fetchWithAuth(`${API_BASE_URL}/dashboard/charts`)
            ]);

            // 1. Set data ringkasan (StatCards)
            // Backend mengirim: { "data": { "total_belanja": ... } }
            if (summaryData && summaryData.data) {
                setSummary(summaryData.data);
            } else {
                setSummary({ total_belanja: 0, budget: 0, sisa_budget: 0 });
            }

            // 2. Set data chart
            // Backend mengirim: { "data": { "pie_chart": [...], "bar_chart": [...] } }
            if (chartsData && chartsData.data) {
                setPieData(chartsData.data.pie_chart || []);
                setBarData(chartsData.data.bar_chart || []);
            } else {
                setPieData([]);
                setBarData([]);
            }
            
        } catch (err) {
            setError(err.message);
            console.error("Gagal mengambil data dashboard:", err);
        } finally {
            setIsLoading(false);
        }
    }, [fetchWithAuth]); // Tambahkan dependensi fetchWithAuth

    // Hook untuk menjalankan fetchData() saat komponen dimuat
    // --- PERBAIKAN 4: Panggil fetchData ---
    useEffect(() => {
        fetchData();
    }, [fetchData]);
    
    // Data dan warna untuk Pie Chart (tidak berubah)
    const PIE_COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#AF19FF'];

    // Format mata uang Rupiah (tidak berubah)
    const formatIDR = (value) => 
        new Intl.NumberFormat('id-ID', { 
            style: 'currency', 
            currency: 'IDR', 
            minimumFractionDigits: 0 
        }).format(value || 0);

    // Tampilan Loading (tidak berubah)
    if (isLoading) {
        return (
            <div className="flex justify-center items-center h-96">
                <Loader className="animate-spin h-12 w-12 text-indigo-600" />
                <p className="ml-4 text-lg text-gray-700">Memuat data dasbor...</p>
            </div>
        );
    }
    
    // Tampilan Error Utama (tidak berubah)
    if (error) {
        return (
            <div className="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded-lg" role="alert">
                <div className="flex">
                    <div className="py-1"><AlertCircle className="h-6 w-6 text-red-500 mr-3" /></div>
                    <div>
                        <p className="font-bold">Gagal Mengambil Data</p>
                        <p className="text-sm">{error}</p>
                    </div>
                </div>
            </div>
        );
    }

    // Tampilan Utama (sesuai mockup TK2 Hal. 21)
    return (
        <div className="space-y-6">
            {/* 1. Kartu Statistik (tidak berubah) */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <StatCard 
                    title="Total Belanja Mingguan" 
                    value={formatIDR(summary.total_belanja)}
                    icon={DollarSign}
                    colorClass="text-red-600"
                    bgColorClass="bg-red-50"
                />
                <StatCard 
                    title="Budget Mingguan" 
                    value={formatIDR(summary.budget)}
                    icon={Target}
                    colorClass="text-blue-600"
                    bgColorClass="bg-blue-50"
                />
                <StatCard 
                    title="Sisa Budget" 
                    value={formatIDR(summary.sisa_budget)}
                    icon={Activity}
                    colorClass={summary.sisa_budget >= 0 ? "text-green-600" : "text-red-600"}
                    bgColorClass={summary.sisa_budget >= 0 ? "bg-green-50" : "bg-red-50"}
                />
            </div>

            {/* 2. Charts */}
            <div className="grid grid-cols-1 lg:grid-cols-5 gap-6">
                
                {/* Pie Chart (tidak berubah) */}
                <div className="lg:col-span-2">
                    <ChartCard title="Pengeluaran per Kategori" icon={TrendingUp}>
                        {pieData.length > 0 ? (
                            <ResponsiveContainer width="100%" height="100%">
                                <PieChart>
                                    <Pie
                                        data={pieData}
                                        cx="50%"
                                        cy="50%"
                                        labelLine={false}
                                        outerRadius={100}
                                        fill="#8884d8"
                                        dataKey="value" // 'value' sesuai dengan model PieChartItem di backend
                                        label={({ name, percent }) => `${(percent * 100).toFixed(0)}%`}
                                    >
                                        {pieData.map((entry, index) => (
                                            <Cell key={`cell-${index}`} fill={PIE_COLORS[index % PIE_COLORS.length]} />
                                        ))}
                                    </Pie>
                                    <Tooltip formatter={(value) => formatIDR(value)} />
                                    <Legend />
                                </PieChart>
                            </ResponsiveContainer>
                        ) : (
                            <div className="flex items-center justify-center h-full text-gray-500">
                                Belum ada data kategori.
                            </div>
                        )}
                    </ChartCard>
                </div>

                {/* Bar Chart (Pengeluaran Mingguan) */}
                <div className="lg:col-span-3">
                    <ChartCard title="Pengeluaran Mingguan" icon={BarIcon}>
                        {barData.length > 0 ? (
                            <ResponsiveContainer width="100%" height="100%">
                                <BarChart data={barData} margin={{ top: 5, right: 20, left: 10, bottom: 5 }}>
                                    <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                                    <XAxis dataKey="name" stroke="#6b7280" />
                                    <YAxis stroke="#6b7280" tickFormatter={(value) => `${formatIDR(value).replace('Rp', '')}`} />
                                    <Tooltip formatter={(value) => formatIDR(value)} />
                                    <Legend />
                                    
                                    {/* --- PERBAIKAN 5: 'Pengeluaran' -> 'pengeluaran' --- */}
                                    {/* Backend mengirim 'pengeluaran' (lowercase) sesuai JSON tag */}
                                    <Bar dataKey="pengeluaran" fill="#4f46e5" radius={[4, 4, 0, 0]} />
                                </BarChart>
                            </ResponsiveContainer>
                        ) : (
                            <div className="flex items-center justify-center h-full text-gray-500">
                                Belum ada data pengeluaran.
                            </div>
                        )}
                    </ChartCard>
                </div>
            </div>
        </div>
    );
};

export default Dashboard;