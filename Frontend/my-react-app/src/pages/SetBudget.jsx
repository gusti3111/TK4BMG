import React, { useState, useEffect, useCallback } from 'react';
import { Loader, AlertCircle, Save, DollarSign, CheckCircle } from 'lucide-react';

const API_BASE_URL = 'http://localhost:8080/api/v1';

/**
 * Halaman Set Budget (sesuai TK2 Hal. 23)
 */
const SetBudget = () => {
    const [nominalBudget, setNominalBudget] = useState('');
    const [summary, setSummary] = useState({
        budget: 0,
        terpakai: 0,
        sisa: 0,
        persen_terpakai: 0
    });
    const [isLoading, setIsLoading] = useState(true);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [error, setError] = useState(null);
    const [submitError, setSubmitError] = useState(null);
    const [submitSuccess, setSubmitSuccess] = useState(false);

    /**
     * Fungsi helper fetchWithAuth (Tidak berubah)
     */
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

    /**
     * ======================================================
     * =============== PERBAIKAN DI FUNGSI INI ==============
     * ======================================================
     * Mengambil data ringkasan budget (progress bar) dari backend
     */
    const fetchBudgetSummary = useCallback(async () => {
        setIsLoading(true);
        setError(null);
        try {
            // PERBAIKAN 1: Endpoint diubah ke /dashboard/summary
            const data = await fetchWithAuth(`${API_BASE_URL}/dashboard/summary`);
            
            // PERBAIKAN 2: Cek 'data' langsung (bukan data.data)
            if (data) {
                // PERBAIKAN 3: Gunakan 'total_belanja' dari backend, bukan 'terpakai'
                const budget = data.budget || 0;
                const terpakai = data.total_belanja || 0; 
                
                // Logika kalkulasi ini sudah benar
                const sisa = budget - terpakai;
                const persen = budget > 0 ? (terpakai / budget) * 100 : 0;
                
                setSummary({
                    budget: summary.budget,
                    terpakai: summary.terpakai,
                    sisa: summary.sisa,
                    persen_terpakai: summary.persen_terpakai
                });
                
                setNominalBudget(budget.toString());
            } else {
                setSummary({ budget: 0, terpakai: 0, sisa: 0, persen_terpakai: 0 });
            }
        } catch (err) {
            setError(err.message);
            console.error("Gagal mengambil summary budget:", err);
        } finally {
            setIsLoading(false);
        }
    }, [fetchWithAuth]);

    // Mengambil data saat komponen dimuat (Tidak berubah)
    useEffect(() => {
        fetchBudgetSummary();
    }, [fetchBudgetSummary]);

    /**
     * Menangani submit form (Set Budget Mingguan)
     * FUNGSI INI SUDAH BENAR - Tidak perlu diubah
     */
    const handleSubmit = async (e) => {
        e.preventDefault();
        
        const amount = parseFloat(nominalBudget);
        if (isNaN(amount) || amount < 0) {
            setSubmitError("Nominal budget tidak valid.");
            return;
        }

        setIsSubmitting(true);
        setSubmitError(null);
        setSubmitSuccess(false);

        try {
            // Endpoint POST /api/v1/budgets sudah benar
            await fetchWithAuth(`${API_BASE_URL}/budgets`, {
                method: 'POST',
                body: JSON.stringify({ jumlah_anggaran: amount })
            });
            
            setSubmitSuccess(true);
            await fetchBudgetSummary(); // Refresh progress bar
            
            setTimeout(() => setSubmitSuccess(false), 3000);

        } catch (err) {
            setSubmitError(err.message);
            console.error("Gagal submit budget:", err);
        } finally {
            setIsSubmitting(false);
        }
    };

    // ... (Sisa kode (formatIDR, JSX/tampilan) tidak perlu diubah) ...
    
    // Format mata uang Rupiah
    const formatIDR = (value) => 
        new Intl.NumberFormat('id-ID', { 
            style: 'currency', 
            currency: 'IDR', 
            minimumFractionDigits: 0 
        }).format(value || 0);

    // Tampilan Loading
    if (isLoading) {
        return (
            <div className="flex justify-center items-center h-64">
                <Loader className="animate-spin h-12 w-12 text-indigo-600" />
            </div>
        );
    }

    // Tampilan Error Utama
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

    // Tampilan Utama (sesuai mockup TK2 Hal. 23)
    return (
        <div className="space-y-6 max-w-4xl mx-auto">
            <h1 className="text-3xl font-extrabold text-gray-900">Set Budget Mingguan</h1>
            <p className="text-gray-600">Atur anggaran belanja mingguan Anda dan pantau penggunaannya. (Sesuai TK2 Hal. 23)</p>

            {/* 1. Form Atur Budget (sesuai mockup) */}
            <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200">
                <h2 className="text-xl font-semibold text-gray-800 mb-4 border-b pb-3">
                    Atur Budget Mingguan
                </h2>
                {/* */}
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label htmlFor="nominalBudget" className="block text-sm font-medium text-gray-700">Nominal Budget (Rp)</label>
                        <div className="mt-1 relative rounded-lg shadow-sm">
                            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <DollarSign className="h-5 w-5 text-gray-400" />
                            </div>
                            <input
                                type="number"
                                id="nominalBudget"
                                value={nominalBudget}
                                onChange={(e) => setNominalBudget(e.target.value)}
                                className="block w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                placeholder="1000000"
                                required
                                min="0"
                            />
                        </div>
                    </div>
                    
                    {submitError && (
                        <div className="bg-red-100 text-red-700 px-4 py-3 rounded-lg" role="alert">
                            {submitError}
                        </div>
                    )}
                    
                    {submitSuccess && (
                        <div className="bg-green-100 text-green-700 px-4 py-3 rounded-lg flex items-center" role="alert">
                            <CheckCircle className="w-5 h-5 mr-2" />
                            Budget berhasil disimpan!
                        </div>
                    )}

                    <button
                        type="submit"
                        disabled={isSubmitting}
                        className="flex items-center justify-center w-full sm:w-auto py-2 px-6 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:bg-gray-400 disabled:cursor-not-allowed transition duration-150"
                    >
                        {isSubmitting ? (
                            <Loader className="animate-spin h-5 w-5 mr-2" /> 
                        ) : (
                            <Save className="h-5 w-5 mr-2" />
                        )}
                        {isSubmitting ? 'Menyimpan...' : 'Simpan Budget'}
                    </button>
                </form>
            </div>

            {/* 2. Progress Penggunaan Budget (sesuai mockup) */}
            <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200">
                <h2 className="text-xl font-semibold text-gray-800 mb-4 border-b pb-3">
                    Progress Penggunaan Budget
                </h2>
                
                {/* Progress Bar */}
                <div className="w-full bg-gray-200 rounded-full h-4 mb-2 overflow-hidden">
                    <div
                        className={`h-4 rounded-full transition-all duration-500 ${summary.persen_terpakai > 85 ? 'bg-red-500' : (summary.persen_terpakai > 60 ? 'bg-yellow-500' : 'bg-green-500')}`}
                        style={{ width: `${summary.persen_terpakai}%` }}
                    ></div>
                </div>
                <p className="text-right text-lg font-bold text-gray-700">
                    {summary.persen_terpakai}% Terpakai
                </p>

                {/* Rincian Teks (sesuai mockup) */}
                <div className="mt-4 space-y-2 text-gray-700">
                    <div className="flex justify-between">
                        <span className="font-medium">Budget Mingguan:</span>
                        <span className="font-bold text-blue-600">{formatIDR(summary.budget)}</span>
                    </div>
                    <div className="flex justify-between">
                        <span className="font-medium">Sudah Dipakai:</span>
                        <span className="font-bold text-red-600">{formatIDR(summary.terpakai)}</span>
                    </div>
                    <hr className="my-1 border-gray-200" />
                    <div className="flex justify-between text-lg">
                        <span className="font-bold">Sisa Budget:</span>
                        <span className={`font-bold ${summary.sisa >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                            {formatIDR(summary.sisa)}
                        </span>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default SetBudget;