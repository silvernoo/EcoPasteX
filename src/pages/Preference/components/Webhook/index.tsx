import { Input } from "antd";
import { useTranslation } from "react-i18next";
import { useSnapshot } from "valtio";
import ProList from "@/components/ProList";
import ProSwitch from "@/components/ProSwitch";
import { globalStore } from "@/stores/global";

const Webhook = () => {
  const { webhook } = useSnapshot(globalStore);
  const { t } = useTranslation();

  return (
    <ProList header={t("preference.settings.webhook_settings.title")}>
      <ProSwitch
        description={t("preference.settings.webhook_settings.hints.enabled")}
        onChange={(value) => {
          globalStore.webhook.enabled = value;
        }}
        title={t("preference.settings.webhook_settings.label.enabled")}
        value={webhook.enabled}
      />

      <div className="flex flex-col gap-2 p-3">
        <div className="font-medium text-color-1 text-sm">
          {t("preference.settings.webhook_settings.label.url")}
        </div>
        <Input
          disabled={!webhook.enabled}
          onChange={(e) => {
            globalStore.webhook.url = e.target.value;
          }}
          placeholder={t(
            "preference.settings.webhook_settings.placeholder.url",
          )}
          value={webhook.url}
        />
        <div className="text-color-3 text-xs">
          {t("preference.settings.webhook_settings.hints.url")}
        </div>
      </div>
    </ProList>
  );
};

export default Webhook;
