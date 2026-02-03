import rss from '@astrojs/rss';

export async function GET(context) {
  // Example blog posts/updates - replace with your actual data source
  const posts = [
    {
      title: 'gsheetbase.com Public Beta Launch',
      description: 'We\'re excited to announce the public beta launch of gsheetbase.com - the fastest way to turn Google Sheets into production-ready JSON APIs.',
      pubDate: new Date('2026-02-01'),
      link: '/blog/public-beta-launch',
    },
    {
      title: 'Type-Safe API Responses with Auto-Generated TypeScript',
      description: 'Learn how gsheetbase automatically generates TypeScript definitions for your spreadsheet data, giving you full type safety and intellisense.',
      pubDate: new Date('2026-01-28'),
      link: '/blog/typescript-definitions',
    },
    {
      title: 'Edge CDN Now Available for All Beta Users',
      description: 'We\'ve rolled out our global CDN to all beta users, ensuring low-latency data delivery no matter where your users are located.',
      pubDate: new Date('2026-01-25'),
      link: '/blog/edge-cdn-rollout',
    },
  ];

  return rss({
    title: 'gsheetbase.com | Updates & News',
    description: 'The latest updates, features, and news from gsheetbase.com - A database that lives in your spreadsheet.',
    site: context.site || 'https://gsheetbase.com',
    items: posts.map((post) => ({
      title: post.title,
      description: post.description,
      pubDate: post.pubDate,
      link: post.link,
    })),
    customData: `<language>en-us</language>`,
    stylesheet: '/rss-styles.xsl', // Optional: add RSS stylesheet
  });
}
