let cachedApiServerUrl = null;

export async function ApiUrl() {
  if (cachedApiServerUrl) {
    return cachedApiServerUrl;
  }

  try {
    const response = await fetch('/meta');
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    const data = await response.json();
    cachedApiServerUrl = data.api_server_url;
    return cachedApiServerUrl;
  } catch (error) {
    console.error('Error fetching API URL:', error);
    return null;
  }
}