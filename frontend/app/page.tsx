'use client';
import React, { useEffect, useState } from 'react';

interface Product {
  id: string;
  name: string;
  price: number;
  description: string;
  image: string;
  stock: number;
}

export default function Home() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [token, setToken] = useState<string | null>(null);
  const [orderMessage, setOrderMessage] = useState<string | null>(null);

  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:9002/ws');

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        setProducts(data);
        setLoading(false);
      } catch (err) {
        console.error('Failed to parse WebSocket message');
      }
    };

    return () => {
      ws.close();
    };
  }, []);

  const handleLogin = async () => {
    try {
      const response = await fetch(`${apiUrl}/api/v1/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: 'user@ecommerce.com', password: '123456' }),
      });
      const result = await response.json();
      if (response.ok) {
        setToken(result.data.token);
        alert('Login Successful!');
      } else {
        alert('Login Failed: ' + result.message);
      }
    } catch (err) {
      alert('Error connecting to Auth Service');
    }
  };

  const handleBuyNow = async (product: Product) => {
    if (!token) {
      alert('Please Login first to place an order!');
      return;
    }

    if (product.stock <= 0) {
      alert('Sorry, this product is out of stock!');
      return;
    }

    try {
      const response = await fetch(`${apiUrl}/api/v1/orders`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ 
          product_id: product.id,
          item: product.name, 
          price: product.price 
        }),
      });
      const result = await response.json();
      if (response.ok) {
        setOrderMessage(`Order Successful! ID: ${result.data.order_id}`);
      } else {
        alert('Order Failed: ' + result.message);
      }
    } catch (err) {
      alert('Error connecting to Order Service');
    }
  };

  return (
    <main className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <header className="flex justify-between items-center bg-white p-6 rounded-xl shadow-xs border border-gray-100">
          <div>
            <h1 className="text-2xl font-bold text-orange-600">Ecommerce WebSocket Real-time Store</h1>
            <p className="text-xs text-gray-500">Instant Stock Synchronization</p>
          </div>
          <div>
            {token ? (
              <span className="bg-emerald-100 text-emerald-800 text-xs px-3 py-1.5 rounded-full font-semibold">
                LoggedIn
              </span>
            ) : (
              <button 
                onClick={handleLogin}
                className="bg-gray-900 text-white px-4 py-2 rounded-lg text-xs font-semibold hover:bg-gray-800 transition"
              >
                Simulate Login
              </button>
            )}
          </div>
        </header>

        {orderMessage && (
          <div className="bg-emerald-50 border border-emerald-200 text-emerald-800 p-4 rounded-xl font-medium text-sm flex justify-between items-center">
            <span>🎉 {orderMessage}</span>
            <button onClick={() => setOrderMessage(null)} className="text-xs underline text-emerald-600">Dismiss</button>
          </div>
        )}

        <section>
          <h2 className="text-xl font-bold text-gray-800 mb-4">Live Product Catalog (WebSocket Connected)</h2>

          {loading && <div className="text-center py-12 text-gray-400">Connecting to WebSocket stream...</div>}

          {!loading && (
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
              {products.map((product) => (
                <div key={product.id} className="bg-white p-4 rounded-xl border border-gray-100 shadow-xs hover:shadow-md transition flex flex-col justify-between">
                  <div>
                    <div className="h-44 bg-gray-50 rounded-lg overflow-hidden flex items-center justify-center mb-3 relative">
                      <img src={product.image} alt={product.name} className="w-full h-full object-cover" />
                      <span className={`absolute top-2 right-2 text-[10px] font-bold px-2 py-1 rounded-md ${product.stock > 0 ? 'bg-emerald-500 text-white' : 'bg-red-500 text-white'}`}>
                        {product.stock > 0 ? `Stock: ${product.stock}` : 'Stock Out'}
                      </span>
                    </div>
                    <h3 className="font-semibold text-gray-800">{product.name}</h3>
                    <p className="text-xs text-gray-500 mt-1">{product.description}</p>
                  </div>

                  <div className="flex justify-between items-center mt-4 pt-3 border-t border-gray-100">
                    <span className="text-lg font-extrabold text-orange-600">৳{product.price.toLocaleString()}</span>
                    <button 
                      onClick={() => handleBuyNow(product)}
                      disabled={product.stock <= 0}
                      className={`px-3 py-1.5 rounded-lg text-xs font-semibold transition ${product.stock > 0 ? 'bg-orange-500 text-white hover:bg-orange-600' : 'bg-gray-200 text-gray-400 cursor-not-allowed'}`}
                    >
                      {product.stock > 0 ? 'Buy Now' : 'Sold Out'}
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>
      </div>
    </main>
  );
}