import React from 'react';
import { Clock } from 'lucide-react';

/**
 * Item tunggal untuk feed aktivitas.
 */
const ActivityItem = ({ icon: Icon, description, time, iconColor }) => {
    return (
        <li className="flex items-start space-x-3 py-3 border-b last:border-b-0">
          {Icon && <Icon className={`w-5 h-5 ${iconColor} flex-shrink-0 mt-1`} />}
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-gray-800 truncate">{description}</p>
            <p className="text-xs text-gray-500 mt-0.5 flex items-center">
              <Clock className="w-3 h-3 mr-1" />
              {time}
            </p>
          </div>
        </li>
    );
};

export default ActivityItem;
