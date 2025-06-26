// File ini sengaja dibuat untuk menguji konfigurasi linter Biome.

const LinterTestPage = () => {
  // ⚠️ PERINGATAN (Tidak akan memblokir commit)
  // Aturan: correctness/noUnusedVariables (level: "warn")
  const myUnusedVariable = "Halo dunia!";

  // ℹ️ INFO (Tidak akan memblokir commit)
  // Aturan: suspicious/noConsoleLog (level: "info")
  console.log("Testing linter hook");

  const items = [
    { id: 1, name: "Buku" },
    { id: 2, name: "Pensil" },
  ];

  return (
    <div>
      <h1>Uji Coba Linter</h1>

      <h2>Daftar yang Salah #1: Menggunakan Index sebagai Key</h2>
      <ul>
        {/* ❌ ERROR (Akan memblokir commit) */}
        {/* Aturan: suspicious/noArrayIndexKey (level: "error") */}
        {items.map((item, index) => (
          <li key={index}>{item.name}</li>
        ))}
      </ul>

      <h2>Daftar yang Salah #2: Tidak ada Key sama sekali</h2>
      <ul>
        {/* ❌ ERROR (Akan memblokir commit) */}
        {/* Aturan: correctness/useJsxKeyInIterable (level: "error") */}
        {items.map((item) => (
          <li>{item.name}</li>
        ))}
      </ul>

      <h2>Daftar yang Benar</h2>
      <ul>
        {/* ✅ BENAR (Tidak akan ada error) */}
        {items.map((item) => (
          <li key={item.id}>{item.name}</li>
        ))}
      </ul>
    </div>
  );
};

export default LinterTestPage;
