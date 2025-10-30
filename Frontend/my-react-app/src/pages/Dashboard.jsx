import React, { useState, useEffect, useCallback } from 'react';
import { 
    PieChart, Pie, Cell, ResponsiveContainer, 
    BarChart, Bar, XAxis, YAxis, Tooltip, Legend, CartesianGrid 
} from 'recharts';
import { Loader, AlertCircle, DollarSign, Target, Activity, TrendingUp, BarChart as BarIcon } from 'lucide-react';

// --- MOCK DATA LOKAL (Digunakan untuk development tanpa backend) ---
const mockSummary = {
    total_belanja: 850000,
    budget: 1200000,
    sisa_budget: 350000,
};

const mockCharts = {
    pie_chart: [
        { name: 'Makanan', value: 350000 },
        { name: 'Transportasi', value: 150000 },
        { name: 'Hiburan', value: 200000 },
        { name: 'Tagihan', value: 150000 },
    ],
    bar_chart: [
        { name: 'Minggu 1', Pengeluaran: 400000 },
        { name: 'Minggu 2', Pengeluaran: 300000 },
        { name: 'Minggu 3', Pengeluaran: 550000 },
        { name: 'Minggu 4', Pengeluaran: 250000 },
    ],
};
// --- AKHIR MOCK DATA LOKAL ---


// API Base URL (Dibiarkan untuk referensi integrasi backend)
const API_BASE_URL = 'http://localhost:8080/api/v1';

// --- KOMPONEN LOKAL (Untuk menghindari error impor) ---

/**
 * Komponen StatCard (Sesuai Mockup TK2 & props di Dashboard.jsx)
 * @param {object} props
 * @param {string} props.title - Judul kartu (misal: "Total Belanja Mingguan")
 * @param {string} props.value - Nilai yang ditampilkan (misal: "Rp 850.000")
 * @param {React.Component} props.icon - Ikon dari Lucide
 * @param {string} [props.colorClass] - Kelas warna Tailwind untuk nilai (misal: "text-red-600")
 * @param {string} [props.bgColorClass] - Kelas warna background (misal: "bg-red-50")
 */
const StatCard = ({ title, value, icon: Icon, colorClass = 'text-gray-900', bgColorClass = 'bg-gray-50' }) => (
    <div className={`p-5 rounded-xl shadow-lg border border-gray-200 ${bgColorClass}`}>
        <div className="flex items-center justify-between mb-2">
            <p className="text-sm font-medium text-gray-500">{title}</p>
            {Icon && <Icon className={`w-6 h-6 ${colorClass.replace('text-', 'text-')}`} />}
        </div>
        <p className={`text-3xl font-bold ${colorClass}`}>{value}</p>
    </div>
);

/**
 * Komponen ChartCard (Wrapper - di-inline dari components/ChartCard.jsx)
 * @param {object} props
 * @param {string} props.title - Judul chart
 * @param {React.Node} props.children - Komponen chart dari Recharts
 * @param {React.Component} [props.icon] - Ikon dari Lucide
 */
const ChartCard = ({ title, children, icon: Icon }) => (
  <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200 h-full flex flex-col">
    <div className="flex items-center justify-between border-b pb-3 mb-4">
        <h3 className="text-xl font-semibold text-gray-800">{title}</h3>
        {Icon && <Icon className="w-6 h-6 text-indigo-500" />}
    </div>
    {/* Div dengan tinggi fleksibel untuk menampung ResponsiveContainer dari Recharts */}
    <div className="flex-grow w-full h-72 min-h-[288px]"> 
      {children}
    </div>
  </div>
);
// --- AKHIR KOMPONEN LOKAL ---


/**
 * Komponen Halaman Dashboard (sesuai TK2 - Hal 21)
 * Menggunakan StatCard dan ChartCard
 */
const Dashboard = () => {
    
    // State untuk data ringkasan (kartu)
    const [summary, setSummary] = useState({
        total_belanja: 0,
        budget: 0,
        sisa_budget: 0,
    });
    
    // State untuk data chart
    const [pieData, setPieData] = useState([]);
    const [barData, setBarData] = useState([]);
    
    // State UI
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);

    /**
     * FUNGSI SIMULASI MENGAMBIL DATA (Menggantikan fetchWithAuth)
     * Ini menggunakan data dari mockData yang di-inline
     */
    const simulateFetchData = useCallback(() => {
        setIsLoading(true);
        setError(null);

        // Simulasi penundaan jaringan 500ms
        setTimeout(() => {
            try {
                // 1. Ambil data ringkasan dari mockData
                setSummary(mockSummary);
                
                // 2. Ambil data chart dari mockData
                setPieData(mockCharts.pie_chart || []);
                setBarData(mockCharts.bar_chart || []);
                
            } catch (e) {
                 setError('Gagal memuat data mock: Format data mungkin salah.');
            } finally {
                setIsLoading(false);
            }
        }, 500); // Penundaan 500ms

    }, []);

    // Hook untuk menjalankan fetchData() saat komponen dimuat
    // Sekarang memanggil fungsi simulasi
    useEffect(() => {
        simulateFetchData();
    }, [simulateFetchData]);
    
    // Data dan warna untuk Pie Chart
    const PIE_COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#AF19FF'];

    // Format mata uang Rupiah
    const formatIDR = (value) => 
        new Intl.NumberFormat('id-ID', { 
            style: 'currency', 
            currency: 'IDR', 
            minimumFractionDigits: 0 
        }).format(value || 0); // Pastikan value tidak null/undefined

    // Tampilan Loading
    if (isLoading) {
        return (
            <div className="flex justify-center items-center h-96">
                <Loader className="animate-spin h-12 w-12 text-indigo-600" />
                <p className="ml-4 text-lg text-gray-700">Memuat data dasbor...</p>
            </div>
        );
    }
    
    // Tampilan Error Utama (Masih dipertahankan, walau tidak akan terpicu jika mock data benar)
    if (error) {
        return (
            <div className="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded-lg" role="alert">
                <p className="font-bold">Gagal Mengambil Data</p>
                <p className="text-sm">{error}</p>
            </div>
        );
    }

    // Tampilan Utama (sesuai mockup TK2 Hal. 21)
    return (
        <div className="space-y-6">
            {/* 1. Kartu Statistik (Sesuai Mockup TK2) */}
            {/* Dibuat responsif: 1 kolom di HP, 3 kolom di tablet/desktop */}
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

            {/* 2. Charts (Sesuai Mockup TK2) */}
            {/* Dibuat responsif: 1 kolom di HP, 2 kolom (berbeda rasio) di layar besar */}
            <div className="grid grid-cols-1 lg:grid-cols-5 gap-6">
                
                {/* Pie Chart (Pengeluaran per Kategori) */}
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
                                        dataKey="value"
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
                                    <Bar dataKey="Pengeluaran" fill="#4f46e5" radius={[4, 4, 0, 0]} />
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
