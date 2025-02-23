import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';

import Heading from '@theme/Heading';
import styles from './index.module.css';
import posthog from 'posthog-js';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link
            className="button button--secondary button--lg"
            to="/docs/intro">
            Docs 📜
          </Link>

          <Link
            className="button button--secondary button--lg"
            to="https://github.com/shinobistack/gokakashi"
            style={{marginLeft: '10px'}}
            >
            GitHub <img src="/img/logo/github.svg" class="github-button-image"/>
          </Link>
        </div>
      </div>
    </header>
  );
}

export default function Home() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={`${siteConfig.title}`}
      description="Description will go into a meta tag in <head />">
      <HomepageHeader />
      <main>
      </main>
    </Layout>
  );
}

posthog.init('phc_yZTtQI33sbvCDSU3mUQY0YsZRwfMeAenGG66HNleRgh',
  {
      api_host: 'https://us.i.posthog.com',
      person_profiles: 'identified_only'
  }
)