import api from '@/utils/axiosInstance';


export const fetchClassification = async ({
  queryKey,
}: {
  queryKey: [string, number, number, string];
}) => {
  const [, page, limit, sort] = queryKey;
  try {
    const res = await api.get(`/classification/logs`, {
      params: {
        page,
        limit,
        sortBy: sort,
      },
    });
    return res.data;
  } catch (error) {
    throw new Error('Failed to fetch classification');
  }
};

export const classifyText = async (text: string) => {
  try {
    const res = await api.post(`/classification`, { text });
    return res.data;
  } catch (error) {
    throw new Error('Classification failed');
  }
};