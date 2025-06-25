import type {ReactNode} from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  description: ReactNode;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Real-time Monitoring',
    description: (
      <>
        Monitor HTTP/HTTPS endpoints and push-based services with real-time status updates.
        Get instant notifications when your services go down with smart alerting to avoid spam.
      </>
    ),
  },
  {
    title: 'Beautiful Status Pages',
    description: (
      <>
        Create public status pages to keep your users informed about service availability.
        Share real-time uptime statistics and maintenance schedules with a professional interface.
      </>
    ),
  },
  {
    title: 'Modern Architecture',
    description: (
      <>
        Built with Go backend and React frontend for performance and reliability.
        Strongly typed, extensible architecture with Docker support for easy deployment.
      </>
    ),
  },
];

function Feature({title, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      {/* <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div> */}
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): ReactNode {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
