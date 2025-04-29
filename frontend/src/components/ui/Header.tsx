import { Flex, Text, Icon, HStack, Image } from "@chakra-ui/react";
import { Link } from "react-router-dom";
import { FiUnlock } from "react-icons/fi";
import ExternalSystemsDialog from "@/components/ui/ExternalSystemsDialog"
export default function Header() {
  return (
    <Flex
      bg="white"
      border="1px solid"
      borderColor="gray.200"
      borderRadius="md"
      px={6}
      py={3}
      align="center"
      justify="space-between"
      shadow="sm"
      my={10}
    >
      <HStack gap={3}>
        <Image src="/logo.png" alt="Logo" boxSize="24px" objectFit="contain" />
        <Text fontWeight="bold" fontSize="lg">
          Prompt Injection Detection
        </Text>
      </HStack>

      <HStack gap={6}>
        <ExternalSystemsDialog/>
        <Link to="/change">
          <HStack gap={1} color="gray.600" _hover={{ color: "blue.500" }}>
            <Icon as={FiUnlock} />
            <Text fontSize="sm">Change Password</Text>
          </HStack>
        </Link>
      </HStack>
    </Flex>
  );
}
