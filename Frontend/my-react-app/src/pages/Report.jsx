import React from 'react';
import { BarChart, Tooltip, ResponsiveContainer, Bar as RechartsBar, XAxis, YAxis, CartesianGrid, Legend } from 'recharts';
import { BarChart3 } from 'lucide-react'; // FIX: BarChart3 sekarang diimpor dari lucide-react
// import ChartCard from '../components/ChartCard'; // Dihapus karena akan didefinisikan di sini
// import { mockReportData } from '../data/mockData'; // Dihapus karena akan didefinisikan di sini

// --- KOMPONEN CHARTCARD (Dipindahkan ke sini) ---
const ChartCard = ({ title, children, icon: Icon }) => (
  <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200 h-full flex flex-col">
    <div className="flex items-center justify-between border-b pb-3 mb-4">
        <h3 className="text-xl font-semibold text-gray-800">{title}</h3>
        {Icon && <Icon className="w-6 h-6 text-indigo-500" />}
    </div>
    <div className="flex-grow w-full h-72"> 
      {children}
    </div>
  </div>
);
// ------------------------------------------------

// --- DATA TIRUAN (Dipindahkan ke sini) ---
const mockReportData = [
    { category: 'Pemasaran', amount: 5000000, percent: 35 },
    { category: 'Pengembangan', amount: 4000000, percent: 28 },
    { category: 'Operasional', amount: 3500000, percent: 24 },
    { category: 'Administrasi', amount: 1800000, percent: 13 },
];
// ------------------------------------------

const BarIcon = BarChart; // Alias untuk menghindari konflik nama

/**
 * Halaman Laporan Pengeluaran
 */
const Reports = () => {
    const totalPengeluaran = mockReportData.reduce((sum, item) => sum + item.amount, 0);

    const formatCurrency = (amount) => {
        return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount);
    };

    return (
        <div className="p-4 sm:p-6 lg:p-8 bg-gray-100 min-h-screen">
            <h2 className="text-3xl font-extrabold text-gray-900 mb-8 flex items-center">
                <BarChart3 className="w-8 h-8 mr-3 text-indigo-600" />
                Laporan Pengeluaran Bulanan
            </h2>

            <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200 mb-8">
                <p className="text-lg font-medium text-gray-500">Total Pengeluaran Bulan Ini</p>
                <p className="text-4xl font-extrabold text-red-600 mt-1">
                    {formatCurrency(totalPengeluaran)}
                </p>
                <p className="text-sm text-gray-500 mt-2">Data per {new Date().toLocaleDateString('id-ID')}</p>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {/* Visualisasi Bar Chart */}
                <div className="lg:col-span-2">
                    <ChartCard title="Pengeluaran Berdasarkan Kategori" icon={BarIcon}>
                        <ResponsiveContainer width="100%" height="100%">
                            <BarChart 
                                data={mockReportData} 
                                layout="vertical" 
                                margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
                            >
                                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                                {/* Menghapus "Rp" dari sumbu X agar lebih bersih */}
                                <XAxis type="number" stroke="#6b7280" tickFormatter={val => formatCurrency(val).replace('Rp', '')} />
                                <YAxis dataKey="category" type="category" stroke="#6b7280" />
                                <Tooltip formatter={(value) => [formatCurrency(value), 'Jumlah']} />
                                <Legend wrapperStyle={{ paddingTop: '10px' }} />
                                <RechartsBar dataKey="amount" fill="#ef4444" name="Jumlah Pengeluaran" radius={[5, 5, 0, 0]} />
                            </BarChart>
                        </ResponsiveContainer>
                    </ChartCard>
                </div>

                {/* Detail Tabel Pengeluaran */}
                <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200">
                    <h3 className="text-xl font-semibold text-gray-800 mb-4 border-b pb-3">Rincian Kategori</h3>
                    <ul className="divide-y divide-gray-100">
                        {mockReportData.map((item, index) => (
                            <li key={index} className="flex justify-between items-center py-3">
                                <div className="flex items-center space-x-3">
                                    <div className={`w-3 h-3 rounded-full ${['bg-red-500', 'bg-blue-500', 'bg-yellow-500', 'bg-green-500'][index % 4]}`}></div>
                                    <p className="font-medium text-gray-700">{item.category}</p>
                                </div>
                                <div className="text-right">
                                    <p className="font-semibold text-gray-900">{formatCurrency(item.amount)}</p>
                                    <p className="text-xs text-gray-500">{item.percent}%</p>
                                </div>
                            </li>
                        ))}
                    </ul>
                </div>
            </div>
        </div>
    );
}

export default Reports;
