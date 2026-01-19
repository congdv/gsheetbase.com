
/**
 * Structured Data (JSON-LD) for SEO
 * Helps search engines understand your website better
 */
export const StructuredData = () => {
  const structuredData = {
    '@context': 'https://schema.org',
    '@type': 'SoftwareApplication',
    name: 'GsheetBase',
    applicationCategory: 'BusinessApplication',
    operatingSystem: 'Web',
    offers: {
      '@type': 'Offer',
      price: '0',
      priceCurrency: 'USD',
    },
    aggregateRating: {
      '@type': 'AggregateRating',
      ratingValue: '5',
      ratingCount: '1',
    },
    description:
      'Transform your Google Sheets into powerful REST APIs in seconds. No coding required. Fast, secure, and scalable API endpoints for your spreadsheet data.',
    url: 'https://gsheetbase.com',
    author: {
      '@type': 'Organization',
      name: 'GSheetBase',
      url: 'https://gsheetbase.com',
    },
    publisher: {
      '@type': 'Organization',
      name: 'GSheetBase',
      logo: {
        '@type': 'ImageObject',
        url: 'https://gsheetbase.com/gsheetbase.svg',
      },
    },
  };

  const organizationData = {
    '@context': 'https://schema.org',
    '@type': 'Organization',
    name: 'Gsheetbase',
    url: 'https://gsheetbase.com',
    logo: 'https://gsheetbase.com/gsheetbase.svg',
    description:
      'Turn Google Sheets into REST APIs instantly',
    sameAs: [
      // Add your social media links here
      // 'https://twitter.com/gsheetbase',
      // 'https://linkedin.com/company/gsheetbase',
    ],
  };

  return (
    <>
      <script type="application/ld+json">
        {JSON.stringify(structuredData)}
      </script>
      <script type="application/ld+json">
        {JSON.stringify(organizationData)}
      </script>
    </>
  );
};

export default StructuredData;
