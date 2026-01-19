import { useEffect } from 'react';

interface SEOProps {
  title?: string;
  description?: string;
  keywords?: string;
  ogImage?: string;
  ogUrl?: string;
  noIndex?: boolean;
}

/**
 * SEO Component for dynamic meta tags
 * Updates document meta tags on route changes
 */
export const SEO: React.FC<SEOProps> = ({
  title,
  description,
  keywords,
  ogImage,
  ogUrl,
  noIndex = false,
}) => {
  useEffect(() => {
    // Update title
    if (title) {
      document.title = `${title} | GSheetBase`;
    }

    // Update or create meta tags
    const updateMetaTag = (name: string, content: string, isProperty = false) => {
      const attribute = isProperty ? 'property' : 'name';
      let element = document.querySelector(`meta[${attribute}="${name}"]`);
      
      if (!element) {
        element = document.createElement('meta');
        element.setAttribute(attribute, name);
        document.head.appendChild(element);
      }
      
      element.setAttribute('content', content);
    };

    // Update description
    if (description) {
      updateMetaTag('description', description);
      updateMetaTag('og:description', description, true);
      updateMetaTag('twitter:description', description);
    }

    // Update keywords
    if (keywords) {
      updateMetaTag('keywords', keywords);
    }

    // Update OG image
    if (ogImage) {
      updateMetaTag('og:image', ogImage, true);
      updateMetaTag('twitter:image', ogImage);
    }

    // Update OG URL
    if (ogUrl) {
      updateMetaTag('og:url', ogUrl, true);
      updateMetaTag('twitter:url', ogUrl);
    }

    // Update OG title
    if (title) {
      updateMetaTag('og:title', `${title} | GSheetBase`, true);
      updateMetaTag('twitter:title', `${title} | GSheetBase`);
    }

    // Update robots
    if (noIndex) {
      updateMetaTag('robots', 'noindex, nofollow');
    }
  }, [title, description, keywords, ogImage, ogUrl, noIndex]);

  return null;
};

export default SEO;
