# 🎯 Production Launch Checklist

## Before Deployment

### URLs & Configuration
- [ ] Update all URLs in `index.html` from `gsheetbase.com` to your domain
- [ ] Update `public/sitemap.xml` with your domain and routes
- [ ] Update `src/components/StructuredData.tsx` with your social links
- [ ] Create `.env.production` with proper API URLs
- [ ] Update `site.webmanifest` theme colors if needed

### Assets
- [ ] Verify all favicon files are in `public/` folder
- [ ] Create 1200x630px OG image (optional but recommended)
- [ ] Test logo (`gsheetbase.svg`) displays correctly
- [ ] Ensure all images are optimized

### Code
- [ ] Run `npm run build` successfully
- [ ] Test build locally with `npm run preview`
- [ ] Fix any TypeScript errors
- [ ] Remove console.logs from production code

## Deployment

### Railway Setup
- [ ] Create new service for `web/ui`
- [ ] Set build command: `npm run build`
- [ ] Set start command: `npm run preview`
- [ ] Add environment variables
- [ ] Configure custom domain
- [ ] Enable HTTPS

### DNS Configuration
- [ ] Point domain to Railway
- [ ] Wait for SSL certificate
- [ ] Verify HTTPS works
- [ ] Test www redirect (if needed)

## Post-Deployment Testing

### Functionality
- [ ] App loads correctly
- [ ] Login with Google works
- [ ] Dashboard accessible
- [ ] API calls work
- [ ] All routes load

### SEO & Meta Tags
- [ ] `robots.txt` accessible: `/robots.txt`
- [ ] `sitemap.xml` accessible: `/sitemap.xml`
- [ ] `site.webmanifest` accessible: `/site.webmanifest`
- [ ] Favicon appears in browser tab
- [ ] View page source - verify meta tags
- [ ] Test with [Open Graph Debugger](https://www.opengraph.xyz/)
- [ ] Test with [Twitter Card Validator](https://cards-dev.twitter.com/validator)
- [ ] Test with [Rich Results Test](https://search.google.com/test/rich-results)

### Social Media Previews
- [ ] Share on Facebook - preview looks good
- [ ] Share on LinkedIn - preview looks good
- [ ] Share on Twitter - preview looks good
- [ ] Share on Slack - preview looks good

### Performance
- [ ] Run Lighthouse audit (Chrome DevTools)
  - [ ] Performance > 90
  - [ ] Accessibility > 90
  - [ ] Best Practices > 90
  - [ ] SEO > 90
- [ ] Test on mobile device
- [ ] Test on different browsers (Chrome, Firefox, Safari)
- [ ] Check loading speed

### PWA
- [ ] Web App Manifest loads correctly
- [ ] Icons display properly
- [ ] "Add to Home Screen" prompt works (mobile)
- [ ] App installable on desktop

## Search Engine Submission

### Google
- [ ] Submit to [Google Search Console](https://search.google.com/search-console)
- [ ] Verify ownership
- [ ] Submit sitemap
- [ ] Request indexing for main pages

### Bing
- [ ] Submit to [Bing Webmaster Tools](https://www.bing.com/webmasters)
- [ ] Verify ownership
- [ ] Submit sitemap

## Optional Enhancements

### Analytics
- [ ] Set up Google Analytics
- [ ] Add analytics to environment variables
- [ ] Verify tracking works

### Monitoring
- [ ] Set up error tracking (Sentry)
- [ ] Set up uptime monitoring
- [ ] Configure alerting

### Security
- [ ] Add security headers
- [ ] Enable CORS properly
- [ ] Review authentication flow
- [ ] Check for exposed secrets

### Marketing
- [ ] Create social media accounts
- [ ] Update social links in StructuredData
- [ ] Create blog/changelog
- [ ] Set up email notifications

## 30 Days After Launch

- [ ] Check Google Search Console for issues
- [ ] Review analytics data
- [ ] Check error logs
- [ ] Monitor performance metrics
- [ ] Gather user feedback
- [ ] Update sitemap if new pages added

## Quarterly Reviews

- [ ] Update dependencies
- [ ] Review and update meta descriptions
- [ ] Check broken links
- [ ] Update OG images if branding changes
- [ ] Review Core Web Vitals
- [ ] Update sitemap

---

## 🎉 Launch Day!

When everything is checked off, you're ready to launch! 🚀

**Pro Tip:** Keep this checklist for future deployments and updates.
