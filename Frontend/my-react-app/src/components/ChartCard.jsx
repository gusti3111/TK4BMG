import React from 'react';

/**
 * Wrapper untuk Chart dan Visualisasi.
 */
const ChartCard = ({ title, children, icon: Icon }) => (
  <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200 h-full flex flex-col">
    <div className="flex items-center justify-between border-b pb-3 mb-4">
        <h3 className="text-xl font-semibold text-gray-800">{title}</h3>
        {Icon && <Icon className="w-6 h-6 text-indigo-500" />}
    </div>
    {/* Div dengan tinggi fleksibel untuk menampung ResponsiveContainer dari Recharts */}
    <div className="flex-grow w-full h-72"> 
      {children}
    </div>
  </div>
);

export default ChartCard;
