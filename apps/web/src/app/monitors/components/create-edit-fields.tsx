import { Card } from "@/components/ui/card";
import { CardContent } from "@/components/ui/card";
import General from "./general";
import Intervals from "./intervals";
import  { Essentials, Other } from "./http";
import Notifications from "./notifications";
import Proxies from "./proxies";

const CreateEditFields = ({
  onNewNotifier,
  onNewProxy
}: {
  onNewNotifier: () => void,
  onNewProxy: () => void
}) => {
  return (
    <>
      <Card>
        <CardContent className="space-y-4">
          <General />
        </CardContent>
      </Card>

      <Card>
        <CardContent className="space-y-4">
          <Essentials />
        </CardContent>
      </Card>

      <Card>
        <CardContent className="space-y-4">
          <Notifications onNewNotifier={onNewNotifier} />
        </CardContent>
      </Card>

      <Card>
        <CardContent className="space-y-4">
          <Proxies onNewProxy={onNewProxy} />
        </CardContent>
      </Card>

      <Card>
        <CardContent className="space-y-4">
          <Intervals />
        </CardContent>
      </Card>


      <Card>
        <CardContent className="space-y-4">
          <Other />
        </CardContent>
      </Card>
    </>
  );
};

export default CreateEditFields;
