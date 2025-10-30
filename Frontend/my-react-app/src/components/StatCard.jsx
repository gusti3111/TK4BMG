import React from 'react';

/**
 * Kartu untuk menampilkan statistik kunci.
 */
const StatCard = ({ icon: Icon, title, value, change, color, bgColor }) => (
  <div className="bg-white p-5 rounded-xl shadow-lg transition duration-300 ease-in-out transform hover:scale-[1.02] border-l-4 border-indigo-500 hover:shadow-xl">
    <div className="flex justify-between items-start">
      <div>
        <p className="text-sm font-medium text-gray-500">{title}</p>
        <p className="text-3xl font-bold text-gray-900 mt-1">{value}</p>
      </div>
      <div className={`p-3 rounded-full ${bgColor} flex-shrink-0`}>
        {/* Pastikan Icon adalah komponen React atau Lucide */}
        {Icon && <Icon className={`w-6 h-6 ${color}`} />}
      </div>
    </div>
    <div className="mt-4 flex items-center text-sm">
      <span className={`font-semibold ${color}`}>
        {change}
      </span>
      <span className="text-gray-500 ml-2">
        vs. bulan lalu
      </span>
    </div>
  </div>
);

export default StatCard;
