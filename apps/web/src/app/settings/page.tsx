import Layout from '@/layout';
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card';
import TimezoneSelector from '../../components/TimezoneSelector';

const SettingsPage = () => {
  return (
    <Layout pageName="Settings">
      <div className="max-w-xl flex flex-col gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Timezone</CardTitle>
            <CardDescription>
              Select your preferred timezone for displaying all times in the app.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <TimezoneSelector />
          </CardContent>
        </Card>
        {/* Add more settings sections here as needed */}
      </div>
    </Layout>
  );
}

export default SettingsPage;
